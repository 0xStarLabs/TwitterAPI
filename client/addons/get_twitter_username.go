package addons

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
	tlsClient "github.com/bogdanfinn/tls-client"
)

// GetTwitterUsername retrieves the username of a Twitter account.
//
// Parameters:
//   - httpClient: the HTTP client to make requests with
//   - cookieClient: manages cookies for the request
//   - config: contains Twitter API configuration and constants
//   - logger: handles logging of operations
//   - csrfToken: CSRF token for request authentication
//
// Returns:
//   - string: the account's username (empty if error)
//   - string: new CSRF token from response
//   - error: any error that occurred, including account status errors
//   - models.ActionStatus: the status of the account
func GetTwitterUsername(httpClient tlsClient.HttpClient, cookieClient *utils.CookieClient, config *models.Config, logger utils.Logger, csrfToken string) (string, string, error, models.ActionStatus) {
	for i := 0; i < config.MaxRetries; i++ {
		if i > 0 { // Don't sleep on first try
			utils.RandomSleep(1, 5)
		}

		// Build URL with query parameters
		baseURL := "https://api.x.com/graphql/UhddhjWCl-JMqeiG4vPtvw/Viewer"
		params := url.Values{}
		params.Add("variables", `{"withCommunitiesMemberships":true}`)
		params.Add("features", `{"rweb_tipjar_consumption_enabled":true,"responsive_web_graphql_exclude_directive_enabled":true,"verified_phone_label_enabled":false,"creator_subscriptions_tweet_preview_api_enabled":true,"responsive_web_graphql_skip_user_profile_image_extensions_enabled":false,"responsive_web_graphql_timeline_navigation_enabled":true}`)
		params.Add("fieldToggles", `{"isDelegate":false,"withAuxiliaryUserLabels":false}`)
		fullURL := baseURL + "?" + params.Encode()

		// Create request config with required headers
		reqConfig := utils.DefaultConfig()
		reqConfig.Method = "GET"
		reqConfig.URL = fullURL
		reqConfig.Headers = append(reqConfig.Headers,
			utils.HeaderPair{Key: "authorization", Value: config.Constants.BearerToken},
			utils.HeaderPair{Key: "cookie", Value: cookieClient.CookiesToHeader()},
			utils.HeaderPair{Key: "origin", Value: "https://twitter.com"},
			utils.HeaderPair{Key: "referer", Value: "https://twitter.com/"},
			utils.HeaderPair{Key: "x-csrf-token", Value: csrfToken},
			utils.HeaderPair{Key: "x-twitter-active-user", Value: "no"},
		)

		bodyBytes, resp, err := utils.MakeRequest(httpClient, reqConfig)
		if err != nil {
			logger.Warning("Unknown | Failed to make get username request: %s", err.Error())
			continue
		}

		// Update cookies from response
		cookieClient.SetCookieFromResponse(resp)

		// Get new CSRF token
		newCsrfToken, ok := cookieClient.GetCookieValue("ct0")
		if !ok {
			logger.Error("Unknown | Failed to get new csrf token")
			continue
		}

		// Parse response and handle different account states
		switch {
		case strings.Contains(string(bodyBytes), "screen_name"):
			var responseData getUsernameJSON
			if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
				logger.Error("Unknown | Failed to unmarshal response: %s", err.Error())
				continue
			}
			username := responseData.Data.Viewer.UserResults.Result.Legacy.ScreenName
			logger.Success("%s | Successfully got username", username)
			return username, newCsrfToken, nil, models.StatusSuccess

		case strings.Contains(string(bodyBytes), "this account is temporarily locked"):
			logger.Error("Unknown | Account is temporarily locked!")
			return "", newCsrfToken, models.ErrAccountLocked, models.StatusLocked

		case strings.Contains(string(bodyBytes), "Could not authenticate you"):
			logger.Error("Unknown | Could not authenticate you")
			return "", newCsrfToken, models.ErrAuthFailed, models.StatusAuthError

		default:
			logger.Error("Unknown | Unknown response: %s", string(bodyBytes))
		}
	}

	logger.Error("Unknown | Unable to get twitter username after %d retries", config.MaxRetries)
	return "", "", models.ErrUnknown, models.StatusUnknown
}

// getUsernameJSON represents the JSON response structure from Twitter's GraphQL API
type getUsernameJSON struct {
	Data struct {
		Viewer struct {
			UserResults struct {
				Result struct {
					Legacy struct {
						ScreenName string `json:"screen_name"`
					} `json:"legacy"`
				} `json:"result"`
			} `json:"user_results"`
		} `json:"viewer"`
	} `json:"data"`
}
