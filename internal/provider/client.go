// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
		err = handleAPIResponse(self.Context, response, err, []int{200})
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

// There is not default in Golang, ... So less is better (not defining delay neither ...)
// https://stackoverflow.com/questions/19612449
func DeleteIt(
	client *resty.Client,
	ctx context.Context,
	name string,
	path string,
) diag.Diagnostics {
	diags := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("Deleting %s...", name))
	response, err := client.R().Delete(path)
	err = handleAPIResponse(ctx, response, err, []int{200, 204})
	if err != nil {
		diags.AddError(
			"Client error",
			fmt.Sprintf("Unable to delete %s, got error: %s", name, err))
		return diags
	}

	for retry := range []int{0, 1, 2, 3, 4} {
		time.Sleep(time.Duration(retry) * time.Second)
		tflog.Debug(ctx, fmt.Sprintf("Poll %d of 5 - Check %s is deleted...", retry+1, name))
		response, err := client.R().Get(path)
		err = handleAPIResponse(ctx, response, err, []int{200, 404})
		if err != nil {
			diags.AddError(
				"Client error",
				fmt.Sprintf("Unable to poll %s will deleting it, got error: %s", name, err))
			return diags
		}
		if response.StatusCode() == 404 {
			tflog.Debug(ctx, fmt.Sprintf("Deleted %s successfully", name))
			return diags
		}
	}
	time.Sleep(5 * time.Second)

	diags.AddError("Client error", fmt.Sprintf("Unable to delete %s, its still available!", name))
	return diags
}

func handleAPIResponse(
	ctx context.Context,
	response *resty.Response,
	err error,
	statusCodes []int,
) error {
	tflog.Debug(ctx, "Response Info:")
	tflog.Debug(ctx, fmt.Sprintf("  Error      : %s", err))
	tflog.Debug(ctx, fmt.Sprintf("  Status Code: %d", response.StatusCode()))
	tflog.Debug(ctx, fmt.Sprintf("  Status     : %s", response.Status()))
	tflog.Debug(ctx, fmt.Sprintf("  Proto      : %s", response.Proto()))
	tflog.Debug(ctx, fmt.Sprintf("  Time       : %s", response.Time()))
	tflog.Debug(ctx, fmt.Sprintf("  Received At: %s", response.ReceivedAt()))
	tflog.Debug(ctx, fmt.Sprintf("  Body       : %s", response.String()))

	if err != nil {
		return err
	}

	if !slices.Contains(statusCodes, response.StatusCode()) {
		// https://stackoverflow.com/questions/39595045/convert-int-array-to-string-separated-by
		var statusCodesString []string
		for _, i := range statusCodes {
			statusCodesString = append(statusCodesString, strconv.Itoa(i))
		}
		return errors.New(
			fmt.Sprintf("API response status code %d (expected %s), Body: %s",
				response.StatusCode(),
				strings.Join(statusCodesString, ", "),
				response.String()))
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
