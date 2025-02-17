// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type AriaClient struct {
	// Add whatever fields, client or connection info, etc. here you would need to setup to
	// communicate with the upstream API. Config holds the common attributes that can be
	// passed to the client on initialization.

	// Host must be a the URL to the base of the API.
	Host string

	RefreshToken string `datapolicy:"token"`
	AccessToken  string `datapolicy:"token"`

	OKAPICallsLogLevel string
	KOAPICallsLogLevel string

	// Transport Layer.
	Insecure bool

	// UserAgent is an optional field that specifies the caller of this request.
	UserAgent string

	Context context.Context

	Client *resty.Client

	// Named read-write mutexes for managing resources
	Mutex *RWMutexKV
}

type AccessTokenResponse struct {
	TokenType string `json:"tokenType"`
	Token     string `json:"token"`
}

func (self *AriaClient) Init() diag.Diagnostics {

	diags := self.CheckConfig()
	if diags.HasError() {
		return diags
	}

	client := resty.New()
	client.SetBaseURL(self.Host)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: self.Insecure})
	if len(self.AccessToken) > 0 {
		client.SetAuthToken(self.AccessToken)
	}
	self.Client = client

	diags.Append(self.GetAccessToken()...)

	self.Mutex = NewRWMutexKV()

	return diags
}

func (self *AriaClient) CheckConfig() diag.Diagnostics {
	diags := diag.Diagnostics{}
	if len(self.Host) == 0 {
		diags.AddError("Missing host", "Host is required to request the API")
	}
	if len(self.RefreshToken) == 0 && len(self.AccessToken) == 0 {
		diags.AddError("Missing token", "Either refresh or access token is required")
	}
	return diags
}

func (self *AriaClient) GetAccessToken() diag.Diagnostics {
	diags := diag.Diagnostics{}

	// Refresh access token if refresh token is set and access token is empty
	if len(self.RefreshToken) > 0 && len(self.AccessToken) == 0 {
		self.Debug("Requesting a new API access token at %s", self.Host)

		var token AccessTokenResponse
		response, err := self.Client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{"refreshToken": self.RefreshToken}).
			SetResult(&token).
			Post("iaas/api/login")
		err = self.HandleAPIResponse(response, err, []int{200})
		if err != nil {
			diags.AddError("Unable to retrieve a valid access token", err.Error())
			return diags
		}

		self.AccessToken = token.Token
	}

	if len(self.AccessToken) == 0 {
		diags.AddError(
			"Empty Access Token",
			"Access Token is empty, will be unable to make API calls")
	}

	return diags
}

// Return a new request insance with apiVersion header set, based on path.
func (self AriaClient) R(path string) *resty.Request {
	return self.Client.R().SetQueryParam("apiVersion", GetVersionFromPath(path))
}

func (self AriaClient) ReadIt(
	instance Model,
	instanceRaw APIModel,
	readPath ...string,
) (bool, *resty.Response, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// Path may be given
	var path string
	if len(readPath) == 0 {
		path = instance.ReadPath()
	} else {
		path = readPath[0]
		if len(readPath) > 1 {
			diags.AddError(
				"Internal error",
				fmt.Sprintf("ReadIt - Must have at most 1 readPath, got %d", len(readPath)))
		}
	}

	if path == "" {
		diags.AddError("Internal error", "ReadIt - Instance cannot be retrieved (empty path).")
		return false, nil, diags
	}

	response, err := self.R(path).SetResult(&instanceRaw).Get(path)
	if response.StatusCode() == 404 {
		self.Debug("%s not found", instance.String())
		return false, response, diags
	}

	err = self.HandleAPIResponse(response, err, []int{200})
	if err != nil {
		diags.AddError(
			"Client error",
			fmt.Sprintf("Unable to read %s, got error: %s", instance.String(), err))
	}

	return true, response, diags
}

func (self AriaClient) DeleteIt(
	instance Model,
	conflictMaxAttemptsOptional ...int,
) diag.Diagnostics {
	// Default value, see https://stackoverflow.com/questions/19612449
	conflictMaxAttempts := 5
	if len(conflictMaxAttemptsOptional) > 0 {
		conflictMaxAttempts = conflictMaxAttemptsOptional[0]
	}

	diags := diag.Diagnostics{}
	name := instance.String()
	self.Debug("Deleting %s...", name)

	for attempt := 0; attempt <= conflictMaxAttempts; attempt++ {

		// Delete the resource
		deletePath := instance.DeletePath()
		response, err := self.R(deletePath).Delete(deletePath)
		err = self.HandleAPIResponse(response, err, []int{200, 204})
		if err != nil {
			// This is potentially an error that will be solved by the deletion of other resources.
			// We can retry the delete operation after some time to converge to desired state.
			if attempt < conflictMaxAttempts && response.StatusCode() == 409 {
				time.Sleep(time.Duration(3) * time.Second) // TODO better with randomness?
				continue
			}
			// Either its not a conflict error either we have made sufficient attempts...
			diags.AddError(
				"Client error",
				fmt.Sprintf("Unable to delete %s, got error: %s", name, err))
			return diags
		}

		// Poll resource until deleted, but if we cant read it...
		readPath := instance.ReadPath()
		if len(readPath) == 0 {
			self.Debug("Deleted %s successfully (without polling)", name)
			return diags
		}

		for retry := range []int{0, 1, 2, 3, 4} {
			time.Sleep(time.Duration(retry) * time.Second)
			self.Debug("Poll %d of 5 - Check %s is deleted...", retry+1, name)

			response, err := self.R(readPath).Get(readPath)
			err = self.HandleAPIResponse(response, err, []int{200, 404})
			if err != nil {
				diags.AddError(
					"Client error",
					fmt.Sprintf("Unable to poll %s will deleting it, got error: %s", name, err))
				return diags
			}

			if response.StatusCode() == 404 {
				self.Debug("Deleted %s successfully", name)
				return diags
			}
		}
	}

	diags.AddError("Client error", fmt.Sprintf("Unable to delete %s, its still available!", name))
	return diags
}

func (self AriaClient) HandleAPIResponse(
	response *resty.Response,
	err error,
	statusCodes []int,
) error {

	// https://stackoverflow.com/questions/39595045/convert-int-array-to-string-separated-by
	var statusCodesString []string
	for _, i := range statusCodes {
		statusCodesString = append(statusCodesString, strconv.Itoa(i))
	}
	statusCodesText := strings.Join(statusCodesString, ", ")

	if err != nil {
		self.LogAPIResponseInfo(response, err, statusCodesText)
		return err
	}

	if !slices.Contains(statusCodes, response.StatusCode()) {
		err = fmt.Errorf(
			"API response status code %d (expected %s), Body: %s",
			response.StatusCode(),
			statusCodesText,
			response.String())
	}

	self.LogAPIResponseInfo(response, err, statusCodesText)
	return err
}

func (self AriaClient) LogAPIResponseInfo(
	response *resty.Response,
	err error,
	statusCodesText string,
) {
	request := response.Request
	requestBody, requestBodyErr := json.MarshalIndent(request.Body, "", "\t")
	if requestBodyErr != nil {
		requestBody = []byte("<body>")
	}

	var responseBody []byte
	if strings.Contains(request.URL, "icon/api/icons") && request.Method == "GET" {
		responseBody = []byte("<THE ICON>")
	} else {
		var responseData any
		responseBody = response.Body()
		if json.Unmarshal(responseBody, &responseData) == nil {
			body, bodyErr := json.MarshalIndent(responseData, "", "\t")
			if bodyErr == nil {
				responseBody = body
			}
		}
	}

	level := self.OKAPICallsLogLevel
	if err != nil {
		level = self.KOAPICallsLogLevel
	}

	self.Log(level, strings.Join([]string{
		"",
		"Request Info:",
		fmt.Sprintf("  URL         : %s", request.URL),
		fmt.Sprintf("  Method      : %s", request.Method),
		fmt.Sprintf("  Body        : %s", requestBody),
		"Response Info:",
		fmt.Sprintf("  Error       : %s", err),
		fmt.Sprintf("  Status Code : %d", response.StatusCode()),
		fmt.Sprintf("  Expected    : %s", statusCodesText),
		fmt.Sprintf("  Status      : %s", response.Status()),
		fmt.Sprintf("  Proto       : %s", response.Proto()),
		fmt.Sprintf("  Time        : %s", response.Time()),
		fmt.Sprintf("  Received At : %s", response.ReceivedAt()),
		fmt.Sprintf("  Body        : %s", responseBody),
	}, "\n"))
}

func GetIdFromLocation(response *resty.Response) (string, error) {
	location, err := url.Parse(response.Header().Get("Location"))
	if err != nil {
		return "", err
	}

	parts := strings.Split(location.Path, "/")
	return parts[len(parts)-1], nil
}

func GetVersionFromPath(path string) string {
	// TODO Take first element of path before /, then map it (faster)
	if strings.HasPrefix(path, "abx") {
		return ABX_API_VERSION
	}
	if strings.HasPrefix(path, "blueprint") {
		return BLUEPRINT_API_VERSION
	}
	if strings.HasPrefix(path, "catalog") {
		return CATALOG_API_VERSION
	}
	if strings.HasPrefix(path, "event-broker") {
		return EVENT_BROKER_API_VERSION
	}
	if strings.HasPrefix(path, "form-service") {
		return FORM_API_VERSION
	}
	if strings.HasPrefix(path, "iaas") {
		return IAAS_API_VERSION
	}
	if strings.HasPrefix(path, "icon") {
		return ICON_API_VERSION
	}
	if strings.HasPrefix(path, "policy") {
		return POLICY_API_VERSION
	}
	if strings.HasPrefix(path, "project-service") {
		return PROJECT_API_VERSION
	}
	if strings.HasPrefix(path, "properties") {
		return BLUEPRINT_API_VERSION
	}
	if strings.HasPrefix(path, "vco") {
		return ORCHESTRATOR_API_VERION
	}
	panic(fmt.Sprintf("GetVersionFromPath Not Implemented for path %s", path))
}
