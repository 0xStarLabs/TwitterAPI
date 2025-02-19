package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
)

// NewAccount creates a new Twitter account instance.
//
// Parameters:
//   - authToken: the auth token for the account (required)
//   - ct0: the x-csrf-token for the account (optional, "" by default)
//   - proxy: the proxy in user:pass@host:port format (optional, "" by default)
//
// Returns:
//   - Account: the configured account instance
//
// Example:
//
//	// Create account with just auth token
//	account := twitter.NewAccount("auth_token_here", "", "")
//
//	// Create account with auth token and proxy
//	account := twitter.NewAccount("auth_token_here", "", "user:pass@host:port")
//
//	// Create account with all parameters
//	account := twitter.NewAccount("auth_token_here", "csrf_token", "user:pass@host:port")
func NewAccount(authToken, ct0, proxy string) *models.Account {
	return &models.Account{
		AuthToken: authToken,
		Ct0:       ct0,
		Proxy:     proxy,
	}
}

// AccountInfo represents detailed information about a Twitter account
type AccountInfo struct {
	Name         string
	Username     string
	CreationDate string
	Suspended    bool
	Protected    bool
	Verified     bool
	FollowedBy   bool
	Followers    int
	IsFollowing  bool
	FriendsCount int
	TweetCount   int
}

// User represents a user entry in the multi-user response
type User struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	ScreenName  string `json:"screen_name"`
	AvatarURL   string `json:"avatar_image_url"`
	IsSuspended bool   `json:"is_suspended"`
	IsVerified  bool   `json:"is_verified"`
	IsProtected bool   `json:"is_protected"`
	IsAuthValid bool   `json:"is_auth_valid"`
}

// MultiUserResponse represents the response from Twitter's multi-user list endpoint
type MultiUserResponse struct {
	Users []User `json:"users"`
}

// IsValid checks if the account is valid and retrieves its status.
//
// Returns:
//   - AccountInfo: containing account details like:
//   - Username and display name
//   - Account status (suspended, protected, verified)
//   - ActionResponse: containing:
//   - Success: true if check was successful
//   - Error: any error that occurred
//   - Status: the status of the action
//
// Example:
//
//	info, resp := twitter.IsValid()
//	if resp.Success {
//	    if info.Suspended {
//	        fmt.Println("Account is suspended")
//	    } else {
//	        fmt.Printf("Account %s is valid\n", info.Username)
//	    }
//	}
func (t *Twitter) IsValid() (*AccountInfo, *models.ActionResponse) {
	baseURL := fmt.Sprintf("https://api.x.com/1.1/account/multi/list.json")
	// Create request config
	reqConfig := utils.DefaultConfig()
	reqConfig.Method = "GET"
	reqConfig.URL = baseURL
	reqConfig.Headers = append(reqConfig.Headers,
		utils.HeaderPair{Key: "accept", Value: "*/*"},
		utils.HeaderPair{Key: "authorization", Value: t.Config.Constants.BearerToken},
		utils.HeaderPair{Key: "content-type", Value: "application/x-www-form-urlencoded"},
		utils.HeaderPair{Key: "cookie", Value: t.Cookies.CookiesToHeader()},
		utils.HeaderPair{Key: "origin", Value: "https://x.com"},
		utils.HeaderPair{Key: "x-csrf-token", Value: t.Account.Ct0},
		utils.HeaderPair{Key: "x-twitter-active-user", Value: "yes"},
		utils.HeaderPair{Key: "x-twitter-auth-type", Value: "OAuth2Session"},
		utils.HeaderPair{Key: "user-agent", Value: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"},	   
	)

	// Make the request
	bodyBytes, resp, err := utils.MakeRequest(t.Client, reqConfig)
	if err != nil {
		t.Logger.Error("%s | Failed to get account info: %v", t.Account.Username, err)
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
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		var response MultiUserResponse
		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			t.Logger.Error("%s | Failed to parse account info response: %v", t.Account.Username, err)
			return nil, &models.ActionResponse{
				Success: false,
				Error:   err,
				Status:  models.StatusUnknown,
			}
		}

		// Find the current user in the response
		var currentUser *User
		for _, user := range response.Users {
			if strings.EqualFold(user.ScreenName, t.Account.Username) {
				currentUser = &user
				break
			}
		}

		if currentUser == nil {
			return nil, &models.ActionResponse{
				Success: false,
				Error:   fmt.Errorf("account not found in response"),
				Status:  models.StatusUnknown,
			}
		}

		// Check if account is valid and not suspended
		if !currentUser.IsAuthValid {
			t.Logger.Error("%s | Account authentication is invalid", t.Account.Username)
			return &AccountInfo{
					Username:  currentUser.ScreenName,
					Suspended: currentUser.IsSuspended,
				}, &models.ActionResponse{
					Success: false,
					Error:   models.ErrAuthFailed,
					Status:  models.StatusAuthError,
				}
		}

		info := &AccountInfo{
			Username:  currentUser.ScreenName,
			Name:      currentUser.Name,
			Suspended: currentUser.IsSuspended,
			Protected: currentUser.IsProtected,
			Verified:  currentUser.IsVerified,
		}

		if info.Suspended {
			t.Logger.Warning("%s | Account is suspended", t.Account.Username)
		} else {
			t.Logger.Success("%s | Account is active and valid", t.Account.Username)
		}

		return info, &models.ActionResponse{
			Success: true,
			Status:  models.StatusSuccess,
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
	case strings.Contains(bodyString, "User has been suspended"):
		t.Logger.Error("%s | Account is suspended", t.Account.Username)
		return &AccountInfo{
				Username:  t.Account.Username,
				Suspended: true,
			}, &models.ActionResponse{
				Success: true,
				Status:  models.StatusSuccess,
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
