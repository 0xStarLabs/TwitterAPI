package client

import (
	"fmt"

	"github.com/0xStarLabs/TwitterAPI/client/addons"
	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
	tlsClient "github.com/bogdanfinn/tls-client"
)

// Twitter represents a Twitter API client instance
type Twitter struct {
	Account *models.Account
	Client  tlsClient.HttpClient
	Logger  utils.Logger
	Config  *models.Config
	Cookies *utils.CookieClient
}

// NewTwitter creates a new Twitter API client instance
func NewTwitter(account *models.Account, config *models.Config) (*Twitter, error) {
	// If no config provided, use default
	if config == nil {
		config = models.NewConfig()
	}

	twitter := &Twitter{
		Account: account,
		Logger:  utils.NewLogger(config.LogLevel),
		Config:  config,
		Cookies: utils.NewCookieClient(),
	}

	// Initialize the client
	if err := twitter.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize Twitter client: %w", err)
	}

	return twitter, nil
}

// init initializes the Twitter client
func (t *Twitter) init() error {
	for i := 0; i < t.Config.MaxRetries; i++ {
		if i > 0 { // Don't sleep on first try
			utils.RandomSleep(1, 5)
		}

		// Create HTTP client
		client, err := utils.CreateHttpClient(t.Account.Proxy)
		if err != nil {
			t.Logger.Error("Failed to create HTTP client: %s", err)
			continue
		}
		t.Client = client

		// Set auth cookies
		authToken, ct0, err := addons.SetAuthCookies(i, t.Cookies, t.Account.AuthToken)
		if err != nil {
			t.Logger.Error("Failed to set auth cookies: %s", err)
			continue
		}
		t.Account.AuthToken = authToken
		t.Account.Ct0 = ct0

		// Get username and verify account
		username, newCsrfToken, err, status := addons.GetTwitterUsername(t.Client, t.Cookies, t.Config, t.Logger, t.Account.Ct0)
		if err != nil {
			switch status {
			case models.StatusLocked:
				return fmt.Errorf("account is locked: %w", err)
			case models.StatusAuthError:
				return fmt.Errorf("authentication failed: %w", err)
			case models.StatusInvalidToken:
				return fmt.Errorf("invalid token: %w", err)
			case models.StatusUnknown:
				t.Logger.Error("Unknown error getting username: %s", err)
				continue
			}
		}

		// Update account info
		t.Account.Username = username
		t.Account.Ct0 = newCsrfToken

		t.Logger.Success("%s | Successfully initialized Twitter client and got username", username)
		return nil
	}

	return fmt.Errorf("failed to initialize after %d retries", t.Config.MaxRetries)
}
