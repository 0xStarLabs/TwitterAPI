package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
)

// Follow follows a user by their username.
//
// Parameters:
//   - username: the Twitter username to follow
//
// Returns an ActionResponse containing:
//   - Success: true if follow was successful
//   - Error: any error that occurred
//   - Status: the status of the action (Success, AuthError, etc.)
//
// Example:
//
//	resp := twitter.Follow("username")
//	if resp.Success {
//	    fmt.Println("Successfully followed user")
//	}
func (t *Twitter) Follow(username string) *models.ActionResponse {
	// Build URL and request body
	baseURL := "https://twitter.com/i/api/1.1/friendships/create.json"
	data := url.Values{}
	data.Set("include_profile_interstitial_type", "1")
	data.Set("include_blocking", "1")
	data.Set("include_blocked_by", "1")
	data.Set("include_followed_by", "1")
	data.Set("include_want_retweets", "1")
	data.Set("include_mute_edge", "1")
	data.Set("include_can_dm", "1")
	data.Set("include_can_media_tag", "1")
	data.Set("skip_status", "1")
	data.Set("screen_name", username)

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
		utils.HeaderPair{Key: "origin", Value: "https://twitter.com"},
		utils.HeaderPair{Key: "referer", Value: fmt.Sprintf("https://twitter.com/%s", username)},
		utils.HeaderPair{Key: "x-csrf-token", Value: t.Account.Ct0},
		utils.HeaderPair{Key: "x-twitter-active-user", Value: "yes"},
		utils.HeaderPair{Key: "x-twitter-auth-type", Value: "OAuth2Session"},
	)

	// Make the request
	bodyBytes, resp, err := utils.MakeRequest(t.Client, reqConfig)
	if err != nil {
		t.Logger.Error("%s | Failed to follow %s: %v", t.Account.Username, username, err)
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
			t.Logger.Error("%s | Failed to parse follow response: %v", t.Account.Username, err)
			return &models.ActionResponse{
				Success: false,
				Error:   err,
				Status:  models.StatusUnknown,
			}
		}

		// Check if we got a valid user response (contains screen_name)
		if response.ScreenName != "" {
			t.Logger.Success("%s | Successfully followed %s", t.Account.Username, username)
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
