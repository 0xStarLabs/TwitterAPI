package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
)

// UserInfoResponse represents the GraphQL response for user info
type UserInfoResponse struct {
	Data struct {
		User struct {
			Result struct {
				RestID string `json:"rest_id"`
				Legacy struct {
					Following            bool   `json:"following"`
					CreatedAt            string `json:"created_at"`
					Description          string `json:"description"`
					FavouritesCount      int    `json:"favourites_count"`
					FollowersCount       int    `json:"followers_count"`
					FriendsCount         int    `json:"friends_count"`
					ListedCount          int    `json:"listed_count"`
					Location             string `json:"location"`
					MediaCount           int    `json:"media_count"`
					Name                 string `json:"name"`
					NormalFollowersCount int    `json:"normal_followers_count"`
					ScreenName           string `json:"screen_name"`
					StatusesCount        int    `json:"statuses_count"`
					Verified             bool   `json:"verified"`
				} `json:"legacy"`
				IsBlueVerified bool `json:"is_blue_verified"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

// GetUserInfoByUsername retrieves detailed information about any Twitter user by their username.
//
// Parameters:
//   - username: the Twitter username to look up
//
// Returns:
//   - UserInfoResponse: containing detailed user information like:
//   - User ID and screen name
//   - Profile information
//   - Account statistics
//   - ActionResponse: containing:
//   - Success: true if lookup was successful
//   - Error: any error that occurred
//   - Status: the status of the action
//
// Example:
//
//	info, resp := twitter.GetUserInfoByUsername("username")
//	if resp.Success {
//	    fmt.Printf("User ID: %s\n", info.Data.User.Result.RestID)
//	    fmt.Printf("Followers: %d\n", info.Data.User.Result.Legacy.FollowersCount)
//	}
func (t *Twitter) GetUserInfoByUsername(username string) (*UserInfoResponse, *models.ActionResponse) {
	// Build URL with query parameters
	baseURL := "https://x.com/i/api/graphql/32pL5BWe9WKeSK1MoPvFQQ/UserByScreenName"
	variables := fmt.Sprintf(`{"screen_name":"%s"}`, username)
	features := `{"hidden_profile_subscriptions_enabled":true, "subscriptions_feature_can_gift_premium": true, "profile_label_improvements_pcf_label_in_post_enabled":true,"rweb_tipjar_consumption_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"subscriptions_verification_info_is_identity_verified_enabled":true,"subscriptions_verification_info_verified_since_enabled":true,"highlights_tweets_tab_ui_enabled":true,"responsive_web_twitter_article_notes_tab_enabled":true,"creator_subscriptions_tweet_preview_api_enabled":true,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true}`
	fieldToggles := `{"withAuxiliaryUserLabels":false}`

	params := url.Values{}
	params.Add("variables", variables)
	params.Add("features", features)
	params.Add("fieldToggles", fieldToggles)
	fullURL := baseURL + "?" + params.Encode()

	// Create request config
	reqConfig := utils.DefaultConfig()
	reqConfig.Method = "GET"
	reqConfig.URL = fullURL
	reqConfig.Headers = append(reqConfig.Headers,
		utils.HeaderPair{Key: "accept", Value: "*/*"},
		utils.HeaderPair{Key: "authorization", Value: t.Config.Constants.BearerToken},
		utils.HeaderPair{Key: "content-type", Value: "application/json"},
		utils.HeaderPair{Key: "cookie", Value: t.Cookies.CookiesToHeader()},
		utils.HeaderPair{Key: "referer", Value: fmt.Sprintf("https://x.com/%s", username)},
		utils.HeaderPair{Key: "x-csrf-token", Value: t.Account.Ct0},
		utils.HeaderPair{Key: "x-twitter-active-user", Value: "no"},
		utils.HeaderPair{Key: "x-twitter-auth-type", Value: "OAuth2Session"},
		utils.HeaderPair{Key: "x-twitter-client-language", Value: "en"},
	)

	// Make the request
	bodyBytes, resp, err := utils.MakeRequest(t.Client, reqConfig)
	if err != nil {
		t.Logger.Error("%s | Failed to get user info for %s: %v", t.Account.Username, username, err)
		return nil, &models.ActionResponse{
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
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 && strings.Contains(bodyString, "screen_name") {
		var response UserInfoResponse
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			t.Logger.Error("%s | Failed to parse user info response: %v", t.Account.Username, err)
			return nil, &models.ActionResponse{
				Success: false,
				Error:   err,
				Status:  models.StatusUnknown,
			}
		}

		if response.Data.User.Result.Legacy.ScreenName != "" {
			t.Logger.Success("%s | Successfully got user info for %s", t.Account.Username, username)
			return &response, &models.ActionResponse{
				Success: true,
				Status:  models.StatusSuccess,
			}
		}
	}

	// Handle error responses
	switch {
	case strings.Contains(bodyString, "this account is temporarily locked"):
		t.Logger.Error("%s | Account is temporarily locked", t.Account.Username)
		return nil, &models.ActionResponse{
			Success: false,
			Error:   models.ErrAccountLocked,
			Status:  models.StatusLocked,
		}
	case strings.Contains(bodyString, "Could not authenticate you"):
		t.Logger.Error("%s | Could not authenticate you", t.Account.Username)
		return nil, &models.ActionResponse{
			Success: false,
			Error:   models.ErrAuthFailed,
			Status:  models.StatusAuthError,
		}
	default:
		t.Logger.Error("%s | Unknown response: %s", t.Account.Username, bodyString)
		return nil, &models.ActionResponse{
			Success: false,
			Error:   fmt.Errorf("unknown response: %s", bodyString),
			Status:  models.StatusUnknown,
		}
	}
}
