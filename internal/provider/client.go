package provider

import (
	"crypto/tls"
	"errors"
	"github.com/go-resty/resty/v2"
	"log"
)

type AriaClientConfig struct {
	// Add whatever fields, client or connection info, etc. here you would need to setup to
	// communicate with the upstream API. Config holds the common attributes that can be
	// passed to the client on initialization.

	// Host must be a the URL to the base of the API.
	Host string

	RefreshToken string `datapolicy:"token"`
	AccessToken string `datapolicy:"token"`

	// Transport Layer.
	Insecure bool

	// UserAgent is an optional field that specifies the caller of this request.
	UserAgent string
}

type AccessTokenResponse struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
}

func (cfg *AriaClientConfig) Check() error {
	if len(cfg.Host) == 0 {
		return errors.New("Host is required to request the API.")
	}
	if len(cfg.RefreshToken) == 0 {
		return errors.New("Refresh token is required to request an access token.")
	}
	return nil
}

func (cfg *AriaClientConfig) Client() *resty.Client {
	client := resty.New()
	client.SetBaseURL(cfg.Host)
	client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: cfg.Insecure })
	if len(cfg.AccessToken) > 0 {
		client.SetAuthToken(cfg.AccessToken)
	}
	return client
}

func (cfg *AriaClientConfig) GetAccessToken() (error) {
	// FIXME Handle refreshing token when required

	// Refresh access token if refresh token is set and access token is empty
	if len(cfg.RefreshToken) > 0 && len(cfg.AccessToken) == 0 {
		if /* logging.IsDebugOrHigher() */ true {
			log.Println("[DEBUG] Requesting a new API access token at", cfg.Host)
		}

		var token AccessTokenResponse
		response, err := cfg.Client().R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{"refreshToken": cfg.RefreshToken}).
			SetResult(&token).
			Post("iaas/api/login")
		err = handleAPIResponse(response, err, 200)
		if err != nil {
			return err
		}

		cfg.AccessToken = token.Token
	}

	if len(cfg.AccessToken) == 0 {
		return errors.New("Access Token cannot be empty")
	}

	return nil
}

func handleAPIResponse(response *resty.Response, err error, statusCode int) error {

	if /* logging.IsDebugOrHigher() && */ (err != nil || response.StatusCode() != statusCode) {
		log.Println("[DEBUG] Response Info:")
		log.Println("[DEBUG]   Error      :", err)
		log.Println("[DEBUG]   Status Code:", response.StatusCode())
		log.Println("[DEBUG]   Status     :", response.Status())
		log.Println("[DEBUG]   Proto      :", response.Proto())
		log.Println("[DEBUG]   Time       :", response.Time())
		log.Println("[DEBUG]   Received At:", response.ReceivedAt())
		log.Println("[DEBUG]   Body       :", response.String())
	}

	if err != nil {
		return err
	}

	if response.StatusCode() != statusCode {
		return errors.New(response.String())
	}

	return nil
}
