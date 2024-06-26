// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type AriaClientConfig struct {
	// Add whatever fields, client or connection info, etc. here you would need to setup to
	// communicate with the upstream API. Config holds the common attributes that can be
	// passed to the client on initialization.

	// Host must be a the URL to the base of the API.
	Host string

	RefreshToken string `datapolicy:"token"`
	AccessToken  string `datapolicy:"token"`

	// Transport Layer.
	Insecure bool

	// UserAgent is an optional field that specifies the caller of this request.
	UserAgent string

	Context context.Context
}

type AccessTokenResponse struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
}

func (self *AriaClientConfig) Check() error {
	if len(self.Host) == 0 {
		return errors.New("Host is required to request the API.")
	}
	if len(self.RefreshToken) == 0 {
		return errors.New("Refresh token is required to request an access token.")
	}
	return nil
}

func (self *AriaClientConfig) Client() *resty.Client {
	client := resty.New()
	client.SetBaseURL(self.Host)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: self.Insecure})
	if len(self.AccessToken) > 0 {
		client.SetAuthToken(self.AccessToken)
	}
	return client
}

func (self *AriaClientConfig) GetAccessToken() error {
	// FIXME Handle refreshing token when required

	// Refresh access token if refresh token is set and access token is empty
	if len(self.RefreshToken) > 0 && len(self.AccessToken) == 0 {
		tflog.Debug(self.Context, "Requesting a new API access token at "+self.Host)

		var token AccessTokenResponse
		response, err := self.Client().R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{"refreshToken": self.RefreshToken}).
			SetResult(&token).
			Post("iaas/api/login")
		err = handleAPIResponse(self.Context, response, err, 200)
		if err != nil {
			return err
		}

		self.AccessToken = token.Token
	}

	if len(self.AccessToken) == 0 {
		return errors.New("Access Token cannot be empty")
	}

	return nil
}

func handleAPIResponse(
	ctx context.Context,
	response *resty.Response,
	err error,
	statusCode int,
) error {
	if err != nil || response.StatusCode() != statusCode {
		tflog.Debug(ctx, "Response Info:")
		tflog.Debug(ctx, fmt.Sprintf("  Error      : %s", err))
		tflog.Debug(ctx, fmt.Sprintf("  Status Code: %d", response.StatusCode()))
		tflog.Debug(ctx, fmt.Sprintf("  Status     : %s", response.Status()))
		tflog.Debug(ctx, fmt.Sprintf("  Proto      : %s", response.Proto()))
		tflog.Debug(ctx, fmt.Sprintf("  Time       : %s", response.Time()))
		tflog.Debug(ctx, fmt.Sprintf("  Received At: %s", response.ReceivedAt()))
		tflog.Debug(ctx, fmt.Sprintf("  Body       : %s", response.String()))
	}

	if err != nil {
		return err
	}

	if response.StatusCode() != statusCode {
		return errors.New(response.String())
	}

	return nil
}

func GetIdFromLocation(response *resty.Response) (string, error) {
	location, err := url.Parse(response.Header().Get("Location"))
	if err != nil {
		return "", err
	}

	parts := strings.Split(location.Path, "/")
	return parts[len(parts)-1], nil
}
