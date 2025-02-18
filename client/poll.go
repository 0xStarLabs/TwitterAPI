package client

import (
	"fmt"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/client/addons"
	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
)

// VotePoll votes in a Twitter poll
// tweetID can be either a tweet URL or tweet ID
// answer is the poll option to vote for
func (t *Twitter) VotePoll(tweetID string, answer string) *models.ActionResponse {
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

	// Get tweet details to extract poll info
	tweetDetails, err := t.getTweetDetails(tweetID)
	if err != nil {
		return &models.ActionResponse{
			Success: false,
			Error:   fmt.Errorf("failed to get tweet details: %w", err),
			Status:  models.StatusUnknown,
		}
	}

	// Extract poll info from tweet details
	pollName := "poll" + strings.Split(strings.Split(tweetDetails, `"name":"poll`)[1], `"`)[0]
	cardID := strings.Split(strings.Split(tweetDetails, `card://`)[1], `"`)[0]

	// Build URL and request body
	baseURL := "https://caps.twitter.com/v2/capi/passthrough/1"
	data := fmt.Sprintf(
		"twitter%%3Astring%%3Acard_uri=card%%3A%%2F%%2F%s&"+
			"twitter%%3Along%%3Aoriginal_tweet_id=%s&"+
			"twitter%%3Astring%%3Aresponse_card_name=%s&"+
			"twitter%%3Astring%%3Acards_platform=Web-12&"+
			"twitter%%3Astring%%3Aselected_choice=%s",
		cardID, tweetID, pollName, answer,
	)

	// Create request config
	reqConfig := utils.DefaultConfig()
	reqConfig.Method = "POST"
	reqConfig.URL = baseURL
	reqConfig.Body = strings.NewReader(data)
	reqConfig.Headers = append(reqConfig.Headers,
		utils.HeaderPair{Key: "accept", Value: "*/*"},
		utils.HeaderPair{Key: "authorization", Value: t.Config.Constants.BearerToken},
		utils.HeaderPair{Key: "content-type", Value: "application/x-www-form-urlencoded"},
		utils.HeaderPair{Key: "cookie", Value: t.Cookies.CookiesToHeader()},
		utils.HeaderPair{Key: "origin", Value: "https://twitter.com"},
		utils.HeaderPair{Key: "referer", Value: "https://twitter.com/"},
		utils.HeaderPair{Key: "x-csrf-token", Value: t.Account.Ct0},
		utils.HeaderPair{Key: "x-twitter-active-user", Value: "yes"},
		utils.HeaderPair{Key: "x-twitter-auth-type", Value: "OAuth2Session"},
	)

	// Make the request
	bodyBytes, resp, err := utils.MakeRequest(t.Client, reqConfig)
	if err != nil {
		t.Logger.Error("%s | Failed to vote in poll: %v", t.Account.Username, err)
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
		t.Logger.Success("%s | Successfully voted in poll %s", t.Account.Username, tweetID)
		return &models.ActionResponse{
			Success: true,
			Status:  models.StatusSuccess,
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

// getTweetDetails gets the details of a tweet, including poll information
func (t *Twitter) getTweetDetails(tweetID string) (string, error) {
	baseURL := fmt.Sprintf(
		"https://twitter.com/i/api/graphql/B9_KmbkLhXt6jRwGjJrweg/TweetDetail?variables="+
			"%%7B%%22focalTweetId%%22%%3A%%22%s%%22%%2C%%22with_rux_injections%%22%%3Afalse%%2C"+
			"%%22includePromotedContent%%22%%3Atrue%%2C%%22withCommunity%%22%%3Atrue%%2C"+
			"%%22withQuickPromoteEligibilityTweetFields%%22%%3Atrue%%2C%%22withBirdwatchNotes%%22%%3Atrue%%2C"+
			"%%22withVoice%%22%%3Atrue%%2C%%22withV2Timeline%%22%%3Atrue%%7D&"+
			"features=%%7B%%22responsive_web_graphql_exclude_directive_enabled%%22%%3Atrue%%2C"+
			"%%22verified_phone_label_enabled%%22%%3Afalse%%2C"+
			"%%22creator_subscriptions_tweet_preview_api_enabled%%22%%3Atrue%%2C"+
			"%%22responsive_web_graphql_timeline_navigation_enabled%%22%%3Atrue%%2C"+
			"%%22responsive_web_graphql_skip_user_profile_image_extensions_enabled%%22%%3Afalse%%2C"+
			"%%22tweetypie_unmention_optimization_enabled%%22%%3Atrue%%2C"+
			"%%22responsive_web_edit_tweet_api_enabled%%22%%3Atrue%%2C"+
			"%%22graphql_is_translatable_rweb_tweet_is_translatable_enabled%%22%%3Atrue%%2C"+
			"%%22view_counts_everywhere_api_enabled%%22%%3Atrue%%2C"+
			"%%22longform_notetweets_consumption_enabled%%22%%3Atrue%%2C"+
			"%%22tweet_awards_web_tipping_enabled%%22%%3Afalse%%2C"+
			"%%22freedom_of_speech_not_reach_fetch_enabled%%22%%3Atrue%%2C"+
			"%%22standardized_nudges_misinfo%%22%%3Atrue%%2C"+
			"%%22tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled%%22%%3Atrue%%2C"+
			"%%22longform_notetweets_rich_text_read_enabled%%22%%3Atrue%%2C"+
			"%%22longform_notetweets_inline_media_enabled%%22%%3Atrue%%2C"+
			"%%22responsive_web_enhance_cards_enabled%%22%%3Afalse%%7D&"+
			"fieldToggles=%%7B%%22withArticleRichContentState%%22%%3Atrue%%7D",
		tweetID,
	)

	// Create request config
	reqConfig := utils.DefaultConfig()
	reqConfig.Method = "GET"
	reqConfig.URL = baseURL
	reqConfig.Headers = append(reqConfig.Headers,
		utils.HeaderPair{Key: "accept", Value: "*/*"},
		utils.HeaderPair{Key: "authorization", Value: t.Config.Constants.BearerToken},
		utils.HeaderPair{Key: "content-type", Value: "application/json"},
		utils.HeaderPair{Key: "cookie", Value: t.Cookies.CookiesToHeader()},
		utils.HeaderPair{Key: "referer", Value: fmt.Sprintf("https://twitter.com/i/status/%s", tweetID)},
		utils.HeaderPair{Key: "x-csrf-token", Value: t.Account.Ct0},
		utils.HeaderPair{Key: "x-twitter-active-user", Value: "yes"},
		utils.HeaderPair{Key: "x-twitter-auth-type", Value: "OAuth2Session"},
		utils.HeaderPair{Key: "x-twitter-client-language", Value: "en"},
	)

	// Make the request
	bodyBytes, resp, err := utils.MakeRequest(t.Client, reqConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get tweet details: %w", err)
	}

	// Update cookies
	t.Cookies.SetCookieFromResponse(resp)
	if newCt0, ok := t.Cookies.GetCookieValue("ct0"); ok {
		t.Account.Ct0 = newCt0
	}

	return string(bodyBytes), nil
}
