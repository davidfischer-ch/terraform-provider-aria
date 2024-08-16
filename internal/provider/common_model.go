// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

type Model interface {
	String() string

	CreatePath() string
	ReadPath() string
	UpdatePath() string
	DeletePath() string
}

type APIModel interface{}
