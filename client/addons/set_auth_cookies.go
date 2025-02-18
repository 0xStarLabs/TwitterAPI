package addons

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/twitter_utils"
	"github.com/0xStarLabs/TwitterAPI/utils"
	http "github.com/bogdanfinn/fhttp"
)

// SetAuthCookies sets authentication cookies for a Twitter client.
// Supports both JSON cookie format and simple auth token.
// 
// Parameters:
//   - accountIndex: index of the account for logging purposes
//   - cookieClient: client's cookie manager
//   - twitterAuth: either JSON cookies or auth token string
//
// Returns:
//   - string: auth token
//   - string: CSRF token
//   - error: any error that occurred
//
// Example:
//
//	// Using auth token
//	authToken, csrfToken, err := SetAuthCookies(0, cookieClient, "auth_token_here")
//	
//	// Using JSON cookies
//	authToken, csrfToken, err := SetAuthCookies(0, cookieClient, "[{\"name\":\"auth_token\",\"value\":\"token\"}]")
func SetAuthCookies(accountIndex int, cookieClient *utils.CookieClient, twitterAuth string) (string, string, error) {
	csrfToken := ""
	authToken := ""
	var err error

	// json cookies
	if strings.Contains(twitterAuth, "[") && strings.Contains(twitterAuth, "]") {
		jsonPart := strings.Split(strings.Split(twitterAuth, "[")[1], "]")[0]
		var cookiesJson []map[string]string
		if err := json.Unmarshal([]byte(jsonPart), &cookiesJson); err != nil {
			return "", "", fmt.Errorf("%d | Failed to decode account json cookies: %v", accountIndex, err)
		}

		for _, cookie := range cookiesJson {
			if name, ok := cookie["name"]; ok {
				value := cookie["value"]
				cookieClient.AddCookies([]http.Cookie{{Name: name, Value: value}})
				if name == "ct0" {
					csrfToken = value
				}
				if name == "auth_token" {
					authToken = value
				}
			}
		}

	// auth token
	} else if len(twitterAuth) < 60 {
		csrfToken, err = twitter_utils.GenerateCSRFToken()
		if err != nil {
			return "", "", fmt.Errorf("%d | Failed to generate CSRF token: %v", accountIndex, err)
		}
		
		cookieClient.AddCookies([]http.Cookie{
			{Name: "auth_token", Value: twitterAuth},
			{Name: "ct0", Value: csrfToken},
			{Name: "des_opt_in", Value: "Y"},
		})
		authToken = twitterAuth
	}

	if csrfToken == "" {
		return "", "", errors.New("failed to get csrf token")
	}

	return authToken, csrfToken, nil
}
