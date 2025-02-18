package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/client/addons"
	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
)

// CommentOptions contains optional parameters for creating a comment.
// Currently supports adding media (images) to comments.
type CommentOptions struct {
	MediaBase64 string // Base64 encoded media (optional)
}

// Comment adds a comment to a tweet with optional media attachment.
// 
// Parameters:
//   - content: the text content of the comment
//   - tweetID: the ID or URL of the tweet to comment on
//   - opts: optional parameters like media (can be nil)
//
// Returns an ActionResponse containing:
//   - Success: true if comment was posted
//   - Error: any error that occurred
//   - Status: the status of the action
//
// Example:
//
//	// Simple comment
//	resp := twitter.Comment("Great tweet!", "1234567890", nil)
//	
//	// Comment with media
//	resp := twitter.Comment("Check this out!", "1234567890", &CommentOptions{
//	    MediaBase64: imageBase64,
//	})
//	
//	if resp.Success {
//	    fmt.Println("Successfully posted comment")
//	}
func (t *Twitter) Comment(content string, tweetID string, opts *CommentOptions) *models.ActionResponse {
	// Extract tweet ID if URL was provided
	if strings.Contains(tweetID, "twitter.com") || strings.Contains(tweetID, "x.com") {
		var err error
		tweetID, err = addons.ExtractTweetID(tweetID, t.Account.Username, t.Logger)
		if err != nil {
			return &models.ActionResponse{
				Success: false,
				Error:   fmt.Errorf("invalid tweet URL: %w", err),
				Status:  models.StatusUnknown,
			}
		}
	}

	// If media is provided, upload it first
	var mediaID string
	if opts != nil && opts.MediaBase64 != "" {
		var err error
		mediaID, err = t.UploadMedia(opts.MediaBase64)
		if err != nil {
			return &models.ActionResponse{
				Success: false,
				Error:   fmt.Errorf("failed to upload media: %w", err),
				Status:  models.StatusUnknown,
			}
		}
	}

	// Build URL and request body
	baseURL := "https://twitter.com/i/api/graphql/" + t.Config.Constants.QueryID.Tweet + "/CreateTweet"

	// Build variables based on options
	variables := map[string]interface{}{
		"tweet_text": content,
		"reply": map[string]interface{}{
			"in_reply_to_tweet_id":   tweetID,
			"exclude_reply_user_ids": []string{},
		},
		"dark_request":            false,
		"semantic_annotation_ids": []string{},
	}

	// Add media if provided
	if mediaID != "" {
		variables["media"] = map[string]interface{}{
			"media_entities": []map[string]interface{}{
				{
					"media_id":     mediaID,
					"tagged_users": []string{},
				},
			},
			"possibly_sensitive": false,
		}
	} else {
		variables["media"] = map[string]interface{}{
			"media_entities":     []string{},
			"possibly_sensitive": false,
		}
	}

	// Build the full request body
	requestBody := map[string]interface{}{
		"variables": variables,
		"features": map[string]interface{}{
			"tweetypie_unmention_optimization_enabled":                                true,
			"responsive_web_edit_tweet_api_enabled":                                   true,
			"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
			"view_counts_everywhere_api_enabled":                                      true,
			"longform_notetweets_consumption_enabled":                                 true,
			"responsive_web_twitter_article_tweet_consumption_enabled":                false,
			"tweet_awards_web_tipping_enabled":                                        false,
			"freedom_of_speech_not_reach_fetch_enabled":                               true,
			"standardized_nudges_misinfo":                                             true,
			"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
			"longform_notetweets_rich_text_read_enabled":                              true,
			"longform_notetweets_inline_media_enabled":                                true,
			"responsive_web_graphql_exclude_directive_enabled":                        true,
			"verified_phone_label_enabled":                                            false,
			"responsive_web_media_download_video_enabled":                             false,
			"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
			"responsive_web_graphql_timeline_navigation_enabled":                      true,
			"c9s_tweet_anatomy_moderator_badge_enabled": true,
			"responsive_web_enhance_cards_enabled":      true,
			"rweb_video_timestamps_enabled":             true,	
		},
		"queryId": t.Config.Constants.QueryID.Tweet,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return &models.ActionResponse{
			Success: false,
			Error:   fmt.Errorf("failed to marshal request body: %w", err),
			Status:  models.StatusUnknown,
		}
	}

	// Create request config
	reqConfig := utils.DefaultConfig()
	reqConfig.Method = "POST"
	reqConfig.URL = baseURL
	reqConfig.Body = strings.NewReader(string(jsonBody))
	reqConfig.Headers = append(reqConfig.Headers,
		utils.HeaderPair{Key: "accept", Value: "*/*"},
		utils.HeaderPair{Key: "authorization", Value: t.Config.Constants.BearerToken},
		utils.HeaderPair{Key: "content-type", Value: "application/json"},
		utils.HeaderPair{Key: "cookie", Value: t.Cookies.CookiesToHeader()},
		utils.HeaderPair{Key: "origin", Value: "https://twitter.com"},
		utils.HeaderPair{Key: "referer", Value: "https://twitter.com/compose/tweet"},
		utils.HeaderPair{Key: "x-csrf-token", Value: t.Account.Ct0},
		utils.HeaderPair{Key: "x-twitter-active-user", Value: "yes"},
		utils.HeaderPair{Key: "x-twitter-auth-type", Value: "OAuth2Session"},
	)

	// Make the request
	bodyBytes, resp, err := utils.MakeRequest(t.Client, reqConfig)
	if err != nil {
		t.Logger.Error("%s | Failed to comment: %v", t.Account.Username, err)
		return &models.ActionResponse{
			Success: false,
			Error:   err,
			Status:  models.StatusUnknown,
		}
	}

	// Update cookies
	t.Cookies.SetCookieFromResponse(resp)
	if newCt0, ok := t.Cookies.GetCookieValue("ct0"); ok {
		t.Account.Ct0 = newCt0
	}

	bodyString := string(bodyBytes)

	// Handle successful responses
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		if strings.Contains(bodyString, "duplicate") {
			t.Logger.Success("%s | Comment was already posted", t.Account.Username)
			return &models.ActionResponse{
				Success: true,
				Status:  models.StatusAlreadyDone,
			}
		}

		var response models.TweetGraphQLResponse
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			t.Logger.Error("%s | Failed to parse comment response: %v", t.Account.Username, err)
			return &models.ActionResponse{
				Success: false,
				Error:   err,
				Status:  models.StatusUnknown,
			}
		}

		if response.Data.CreateTweet.TweetResults.Result.RestID != "" {
			t.Logger.Success("%s | Successfully posted comment", t.Account.Username)
			return &models.ActionResponse{
				Success: true,
				Status:  models.StatusSuccess,
			}
		}
	}

	// Handle error responses
	switch {
	case strings.Contains(bodyString, "this account is temporarily locked"):
		t.Logger.Error("%s | Account is temporarily locked", t.Account.Username)
		return &models.ActionResponse{
			Success: false,
			Error:   models.ErrAccountLocked,
			Status:  models.StatusLocked,
		}
	case strings.Contains(bodyString, "Could not authenticate you"):
		t.Logger.Error("%s | Could not authenticate you", t.Account.Username)
		return &models.ActionResponse{
			Success: false,
			Error:   models.ErrAuthFailed,
			Status:  models.StatusAuthError,
		}
	default:
		t.Logger.Error("%s | Unknown response: %s", t.Account.Username, bodyString)
		return &models.ActionResponse{
			Success: false,
			Error:   fmt.Errorf("unknown response: %s", bodyString),
			Status:  models.StatusUnknown,
		}
	}
}
