// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

const ABX_API_VERSION = "2019-09-12"
const BLUEPRINT_API_VERSION = "2019-09-12"
const CATALOG_API_VERSION = "2020-08-25"

// TODO then ensure its used (check related TODOs).
// 7.6 ?? https://developer.vmware.com/apis/576/#api
const EVENT_BROKER_API_VERSION = ""

const FORM_API_VERSION = "1.0"
const IAAS_API_VERSION = "2021-07-15"

// TODO then ensure its used (check related TODOs).
const ICON_API_VERSION = ""

const ORCHESTRATOR_API_VERION = ""

const POLICY_API_VERSION = "2020-08-25"

const PROJECT_API_VERSION = "2019-01-15"

// TODO then ensure its used (check related TODOs).
const PLATFORM_API_VERSION = ""

// Helpers for documenting attributes in schema ----------------------------------------------------

const IMMUTABLE = " (force recreation on change)"

const JSON_INSTEAD_OF_DYNAMIC_DISCLAIMER = " " +
	"(JSON encoded)\n" +
	"\n" +
	"We should have implemented this attribute as a dynamic type (and not JSON).\n" +
	"Unfortunately Terraform SDK returns this issue:\n" +
	"Dynamic types inside of collections are not currently supported in " +
	"terraform-plugin-framework.\n"
