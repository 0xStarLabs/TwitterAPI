package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
)

// Unfollow unfollows a user by their user ID or username.
//
// Parameters:
//   - userIDOrUsername: can be either a numeric user ID or a Twitter username
//
// Returns an ActionResponse containing:
//   - Success: true if unfollow was successful
//   - Error: any error that occurred
//   - Status: the status of the action (Success, AuthError, etc.)
//
// Example:
//
//	// Unfollow by username
//	resp := twitter.Unfollow("username")
//
//	// Unfollow by ID
//	resp := twitter.Unfollow("1234567890")
//
//	if resp.Success {
//	    fmt.Println("Successfully unfollowed user")
//	}
func (t *Twitter) Unfollow(userIDOrUsername string) *models.ActionResponse {
	// Check if the input is not a numeric ID
	if !utils.IsNumeric(userIDOrUsername) {
		// Get user info to get the numeric ID
		info, resp := t.GetUserInfoByUsername(userIDOrUsername)
		if !resp.Success {
			return resp
		}
		userIDOrUsername = info.Data.User.Result.RestID
	}

	// Build URL and request body
	baseURL := "https://x.com/i/api/1.1/friendships/destroy.json"
	data := url.Values{}
	data.Set("include_profile_interstitial_type", "1")
	data.Set("include_blocking", "1")
	data.Set("include_blocked_by", "1")
	data.Set("include_followed_by", "1")
	data.Set("include_want_retweets", "1")
	data.Set("include_mute_edge", "1")
	data.Set("include_can_dm", "1")
	data.Set("include_can_media_tag", "1")
	data.Set("include_ext_is_blue_verified", "1")
	data.Set("include_ext_verified_type", "1")
	data.Set("include_ext_profile_image_shape", "1")
	data.Set("skip_status", "1")
	data.Set("user_id", userIDOrUsername)

	// Create request config
	reqConfig := utils.DefaultConfig()
	reqConfig.Method = "POST"
	reqConfig.URL = baseURL
	reqConfig.Body = strings.NewReader(data.Encode())
	reqConfig.Headers = append(reqConfig.Headers,
		utils.HeaderPair{Key: "accept", Value: "*/*"},
		utils.HeaderPair{Key: "authorization", Value: t.Config.Constants.BearerToken},
		utils.HeaderPair{Key: "content-type", Value: "application/x-www-form-urlencoded"},
		utils.HeaderPair{Key: "cookie", Value: t.Cookies.CookiesToHeader()},
		utils.HeaderPair{Key: "origin", Value: "https://x.com"},
		utils.HeaderPair{Key: "x-csrf-token", Value: t.Account.Ct0},
		utils.HeaderPair{Key: "x-twitter-active-user", Value: "yes"},
		utils.HeaderPair{Key: "x-twitter-auth-type", Value: "OAuth2Session"},
	)

	// Make the request
	bodyBytes, resp, err := utils.MakeRequest(t.Client, reqConfig)
	if err != nil {
		t.Logger.Error("%s | Failed to unfollow user %s: %v", t.Account.Username, userIDOrUsername, err)
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
		var response models.UserResponse
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			t.Logger.Error("%s | Failed to parse unfollow response: %v", t.Account.Username, err)
			return &models.ActionResponse{
				Success: false,
				Error:   err,
				Status:  models.StatusUnknown,
			}
		}

		// Check if we got a valid user response (contains screen_name)
		if response.ScreenName != "" {
			t.Logger.Success("%s | Successfully unfollowed user %s", t.Account.Username, userIDOrUsername)
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
