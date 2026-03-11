// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

// cleanup sweeps all Aria/vRO test resources identified by the ARIA_PROVIDER_TEST
// and Custom.AriaProviderTest naming conventions.
//
// Usage:
//
//	cleanup [-dry-run] [-force]
//
// Required environment variables:
//
//	ARIA_HOST            Base URL of the Aria/vRO instance (e.g. https://my-aria.example.com)
//	ARIA_REFRESH_TOKEN   Refresh token (mutually exclusive with ARIA_ACCESS_TOKEN)
//	ARIA_ACCESS_TOKEN    Access token  (mutually exclusive with ARIA_REFRESH_TOKEN)
//
// Optional environment variables:
//
//	ARIA_INSECURE                  Set to "true" to skip TLS certificate verification
//	TF_VAR_test_project_id         Project ID used for ABX actions and project-scoped catalog sources
//	TF_VAR_test_catalog_item_id    Catalog item ID used to look up custom forms
//	TF_VAR_test_catalog_item_type  Catalog item type used to look up custom forms
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/davidfischer-ch/terraform-provider-aria/internal/provider"
)

// stdLogger wraps the standard log package so it satisfies provider.CleanupLogger.
type stdLogger struct{}

func (stdLogger) Logf(format string, args ...any) {
	log.Printf(format, args...)
}

func main() {
	dryRun := flag.Bool("dry-run", false, "Print what would be deleted without touching the API")
	force := flag.Bool("force", false, "Bypass dependency checks (?force=true) and usage locks (?ignoreUsage=true)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: cleanup [-dry-run] [-force]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nRequired environment variables:\n")
		fmt.Fprintf(os.Stderr, "  ARIA_HOST            Base URL of the Aria/vRO instance\n")
		fmt.Fprintf(os.Stderr, "  ARIA_REFRESH_TOKEN   Refresh token (mutually exclusive with ARIA_ACCESS_TOKEN)\n")
		fmt.Fprintf(os.Stderr, "  ARIA_ACCESS_TOKEN    Access token  (mutually exclusive with ARIA_REFRESH_TOKEN)\n")
		fmt.Fprintf(os.Stderr, "\nOptional environment variables:\n")
		fmt.Fprintf(os.Stderr, "  ARIA_INSECURE                  Skip TLS certificate verification (\"true\")\n")
		fmt.Fprintf(os.Stderr, "  TF_VAR_test_project_id         Project ID for ABX actions and project-scoped catalog sources\n")
		fmt.Fprintf(os.Stderr, "  TF_VAR_test_catalog_item_id    Catalog item ID for custom forms\n")
		fmt.Fprintf(os.Stderr, "  TF_VAR_test_catalog_item_type  Catalog item type for custom forms\n")
	}

	flag.Parse()

	host := os.Getenv("ARIA_HOST")
	refreshToken := os.Getenv("ARIA_REFRESH_TOKEN")
	accessToken := os.Getenv("ARIA_ACCESS_TOKEN")
	insecure := strings.EqualFold(os.Getenv("ARIA_INSECURE"), "true")

	if host == "" {
		fmt.Fprintln(os.Stderr, "Error: ARIA_HOST is required")
		os.Exit(1)
	}
	if refreshToken == "" && accessToken == "" {
		fmt.Fprintln(os.Stderr, "Error: ARIA_REFRESH_TOKEN or ARIA_ACCESS_TOKEN is required")
		os.Exit(1)
	}

	projectID := os.Getenv("TF_VAR_test_project_id")
	catalogItemID := os.Getenv("TF_VAR_test_catalog_item_id")
	catalogItemType := os.Getenv("TF_VAR_test_catalog_item_type")

	client := &provider.AriaClient{
		Host:               host,
		RefreshToken:       refreshToken,
		AccessToken:        accessToken,
		Insecure:           insecure,
		OKAPICallsLogLevel: "DEBUG",
		KOAPICallsLogLevel: "WARN",
		Context:            context.Background(),
	}

	if diags := client.Init(); diags.HasError() {
		for _, d := range diags.Errors() {
			fmt.Fprintf(os.Stderr, "Error: %s — %s\n", d.Summary(), d.Detail())
		}
		os.Exit(1)
	}

	runner := &provider.CleanupRunner{
		Client: client,
		Log:    stdLogger{},
		DryRun: *dryRun,
		Force:  *force,
	}

	if *dryRun {
		log.Println("Dry-run mode: no resources will be deleted")
	}
	if *force {
		log.Println("Force mode: dependency checks and usage locks will be bypassed")
	}

	// Sweep in dependency order: leaves first, parents last.

	// --- Orchestrator (vRO) ---
	runner.OrchestratorTasks()
	runner.OrchestratorActions()
	runner.OrchestratorWorkflows()
	runner.OrchestratorConfigurations()
	runner.OrchestratorEnvironments()
	runner.OrchestratorEnvironmentRepositories()
	runner.OrchestratorCategories()

	// --- Catalog / Service Broker ---
	if catalogItemID != "" && catalogItemType != "" {
		runner.CustomForms(catalogItemID, catalogItemType)
	}
	if projectID != "" {
		runner.CatalogSourcesInProject(projectID)
	}
	runner.CatalogSources()

	// --- ABX ---
	if projectID != "" {
		runner.ABXActions(projectID)
		runner.CustomResourceABXActions(projectID)
	}
	runner.ABXConstants()

	// --- IaaS / Blueprint ---
	runner.CloudTemplates()
	runner.CustomNamings()
	runner.Tags()

	// --- Custom resources / forms ---
	runner.CustomResources()

	// --- Governance ---
	runner.PropertyGroups()
	runner.Policies()
	runner.Subscriptions()

	// --- Projects (last — may own other resources) ---
	runner.Projects()
}
