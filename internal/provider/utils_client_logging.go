// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (self AriaClient) Error(message string, args ...any) {
	tflog.Error(self.Context, fmt.Sprintf(message, args...))
}

func (self AriaClient) Warn(message string, args ...any) {
	tflog.Warn(self.Context, fmt.Sprintf(message, args...))
}

func (self AriaClient) Debug(message string, args ...any) {
	tflog.Debug(self.Context, fmt.Sprintf(message, args...))
}

func (self AriaClient) Info(message string, args ...any) {
	tflog.Debug(self.Context, fmt.Sprintf(message, args...))
}

func (self AriaClient) Trace(message string, args ...any) {
	tflog.Trace(self.Context, fmt.Sprintf(message, args...))
}

func (self AriaClient) Log(level string, message string, args ...any) {
	// Sorted by occurences to optimize branching a little bit
	if level == "DEBUG" {
		self.Debug(message, args...)
	} else if level == "INFO" {
		self.Info(message, args...)
	} else if level == "TRACE" {
		self.Trace(message, args...)
	} else if level == "WARN" {
		self.Warn(message, args...)
	} else if level == "ERROR" {
		self.Error(message, args...)
	} else {
		self.Debug("Unknown log level %s, defaulting to TRACE", level)
		self.Trace(message, args...)
	}
}
