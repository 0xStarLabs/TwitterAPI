package client

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/0xStarLabs/TwitterAPI/models"
	"github.com/0xStarLabs/TwitterAPI/utils"
)

// UploadMedia uploads media to Twitter and returns the media ID.
// Used internally by Tweet and Comment functions when including media.
// 
// Parameters:
//   - mediaBase64: the base64-encoded image data
//
// Returns:
//   - string: the media ID if successful
//   - error: any error that occurred
//
// Example:
//
//	mediaID, err := twitter.UploadMedia(imageBase64)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use mediaID in Tweet or Comment function
func (t *Twitter) UploadMedia(mediaBase64 string) (string, error) {
	mediaURL := "https://upload.twitter.com/1.1/media/upload.json"
	data := url.Values{}
	data.Set("media_data", mediaBase64)

	reqConfig := utils.DefaultConfig()
	reqConfig.Method = "POST"
	reqConfig.URL = mediaURL
	reqConfig.Body = strings.NewReader(data.Encode())
	reqConfig.Headers = append(reqConfig.Headers,
		utils.HeaderPair{Key: "authorization", Value: t.Config.Constants.BearerToken},
		utils.HeaderPair{Key: "content-type", Value: "application/x-www-form-urlencoded"},
		utils.HeaderPair{Key: "cookie", Value: t.Cookies.CookiesToHeader()},
		utils.HeaderPair{Key: "x-csrf-token", Value: t.Account.Ct0},
	)

	bodyBytes, resp, err := utils.MakeRequest(t.Client, reqConfig)
	if err != nil {
		t.Logger.Error("%s | Failed to upload media: %v", t.Account.Username, err)
		return "", err
	}

	// Update cookies
	t.Cookies.SetCookieFromResponse(resp)
	if newCt0, ok := t.Cookies.GetCookieValue("ct0"); ok {
		t.Account.Ct0 = newCt0
	}

	var response models.MediaUploadResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		t.Logger.Error("%s | Failed to parse media upload response: %v", t.Account.Username, err)
		return "", err
	}

	if response.MediaIDString == "" {
		t.Logger.Error("%s | No media ID in response", t.Account.Username)
		return "", err
	}

	t.Logger.Success("%s | Successfully uploaded media", t.Account.Username)
	return response.MediaIDString, nil
} 