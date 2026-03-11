// Copyright (c) State of Geneva (Switzerland)
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"strings"
)

// Naming conventions used to identify test resources.
const (
	TestPrefix               = "ARIA_PROVIDER_TEST"
	TestCustomResourcePrefix = "Custom.AriaProviderTest"
)

// CleanupLogger is a simple logger interface so CleanupRunner can work both in tests and in the
// standalone cleanup binary without importing testing.T.
type CleanupLogger interface {
	Logf(format string, args ...any)
}

// CleanupRunner holds a client, logger, and operation flags.
type CleanupRunner struct {
	Client *AriaClient
	Log    CleanupLogger
	DryRun bool
	// Force, when true, appends ?force=true (vRO) or ?ignoreUsage=true (tags) to bypass
	// dependency checks and usage locks. Without this flag those deletions are skipped.
	Force bool
}

// cleanupEntry describes a single resource deletion.
type cleanupEntry struct {
	label      string
	deletePath string
}

// vROLinkAttributeAPIModel is a single attribute in a vRO link.
type vROLinkAttributeAPIModel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// vROLinkAPIModel is a single link in a vRO list response.
type vROLinkAPIModel struct {
	Href       string                     `json:"href"`
	Rel        string                     `json:"rel"`
	Attributes []vROLinkAttributeAPIModel `json:"attributes"`
}

// vROLinksAPIModel is the vRO list API response envelope.
type vROLinksAPIModel struct {
	Link  []vROLinkAPIModel `json:"link"`
	Total int               `json:"total"`
}

// contentListAPIModel is a generic paginated list response (iaas, catalog, abx, …).
type contentListAPIModel struct {
	Content []map[string]any `json:"content"`
}

// idFromHref extracts the resource ID from a vRO href such as
// "https://host/vco/api/workflows/some-uuid" → "some-uuid".
func idFromHref(href string) string {
	href = strings.TrimRight(href, "/")
	parts := strings.Split(href, "/")
	return parts[len(parts)-1]
}

// applyCleanups executes (or dry-runs) the provided list of deletions.
func (r *CleanupRunner) applyCleanups(entries []cleanupEntry) {
	for _, e := range entries {
		if r.DryRun {
			r.Log.Logf("Would delete %s", e.label)
			continue
		}
		r.Log.Logf("Cleaning up %s", e.label)
		resp, err := r.Client.R(e.deletePath).Delete(e.deletePath)
		if err := r.Client.HandleAPIResponse(resp, err, []int{200, 204, 404}); err != nil {
			r.Log.Logf("Warning: cannot delete %s: %v", e.label, err)
		}
	}
}

// vROCleanupsByPrefix lists a vRO resource collection and returns entries whose nameField
// attribute starts with TestPrefix.
func (r *CleanupRunner) vROCleanupsByPrefix(
	listPath, label, nameField string,
) []cleanupEntry {
	var raw vROLinksAPIModel
	resp, err := r.Client.R(listPath).SetResult(&raw).Get(listPath)
	if err := r.Client.HandleAPIResponse(resp, err, []int{200}); err != nil {
		r.Log.Logf("Warning: cannot list %s resources: %v", label, err)
		return nil
	}

	var entries []cleanupEntry
	for _, link := range raw.Link {
		var nameValue string
		for _, attr := range link.Attributes {
			if attr.Name == nameField {
				nameValue = attr.Value
				break
			}
		}
		if !strings.HasPrefix(nameValue, TestPrefix) {
			continue
		}
		id := idFromHref(link.Href)
		entries = append(entries, cleanupEntry{
			label:      fmt.Sprintf("%s %q", label, nameValue),
			deletePath: listPath + "/" + id,
		})
	}
	return entries
}

// contentCleanupsByPrefix lists a content-based resource collection and returns entries whose
// nameField starts with TestPrefix.
func (r *CleanupRunner) contentCleanupsByPrefix(
	listPath, nameField string,
) []cleanupEntry {
	var raw contentListAPIModel
	resp, err := r.Client.R(listPath).
		SetQueryParam("size", "10000").
		SetResult(&raw).
		Get(listPath)
	if err := r.Client.HandleAPIResponse(resp, err, []int{200}); err != nil {
		r.Log.Logf("Warning: cannot list resources at %s: %v", listPath, err)
		return nil
	}

	var entries []cleanupEntry
	for _, item := range raw.Content {
		id, _ := item["id"].(string)
		nameRaw, _ := item[nameField].(string)
		if !strings.HasPrefix(nameRaw, TestPrefix) || id == "" {
			continue
		}
		entries = append(entries, cleanupEntry{
			label:      fmt.Sprintf("%s %q", listPath, nameRaw),
			deletePath: listPath + "/" + id,
		})
	}
	return entries
}

// contentCleanupsByPrefixInProject is like contentCleanupsByPrefix but filters by projectId.
func (r *CleanupRunner) contentCleanupsByPrefixInProject(
	listPath, prefix, projectID string,
) []cleanupEntry {
	var raw contentListAPIModel
	resp, err := r.Client.R(listPath).
		SetQueryParam("size", "10000").
		SetQueryParam("projectId", projectID).
		SetResult(&raw).
		Get(listPath)
	if err := r.Client.HandleAPIResponse(resp, err, []int{200}); err != nil {
		r.Log.Logf("Warning: cannot list resources at %s (project %s): %v", listPath, projectID, err)
		return nil
	}

	var entries []cleanupEntry
	for _, item := range raw.Content {
		id, _ := item["id"].(string)
		nameRaw, _ := item["name"].(string)
		if !strings.HasPrefix(nameRaw, prefix) || id == "" {
			continue
		}
		entries = append(entries, cleanupEntry{
			label:      fmt.Sprintf("%s %q (project %s)", listPath, nameRaw, projectID),
			deletePath: listPath + "/" + id,
		})
	}
	return entries
}

// ---------- Public sweep methods -----------------------------------------------------------------

// OrchestratorCategories deletes vRO categories whose name starts with TestPrefix.
func (r *CleanupRunner) OrchestratorCategories() {
	r.applyCleanups(
		r.vROCleanupsByPrefix("vco/api/categories", "Category", "name"),
	)
}

// OrchestratorActions deletes vRO actions whose fqn starts with TestPrefix.
// When Force is true, ?force=true is appended to bypass dependency checks.
func (r *CleanupRunner) OrchestratorActions() {
	entries := r.vROCleanupsByPrefix("vco/api/actions", "Action", "fqn")
	if r.Force {
		for i := range entries {
			entries[i].deletePath += "?force=true"
		}
	}
	r.applyCleanups(entries)
}

// OrchestratorWorkflows deletes vRO workflows whose name starts with TestPrefix.
// When Force is true, ?force=true is appended to bypass dependency checks.
func (r *CleanupRunner) OrchestratorWorkflows() {
	entries := r.vROCleanupsByPrefix("vco/api/workflows", "Workflow", "name")
	if r.Force {
		for i := range entries {
			entries[i].deletePath += "?force=true"
		}
	}
	r.applyCleanups(entries)
}

// OrchestratorConfigurations deletes vRO configurations whose name starts with TestPrefix.
func (r *CleanupRunner) OrchestratorConfigurations() {
	r.applyCleanups(
		r.vROCleanupsByPrefix(
			"vco/api/configurations", "Configuration", "name"),
	)
}

// OrchestratorTasks deletes vRO tasks whose name starts with TestPrefix.
func (r *CleanupRunner) OrchestratorTasks() {
	r.applyCleanups(
		r.vROCleanupsByPrefix("vco/api/tasks", "Task", "name"),
	)
}

// OrchestratorEnvironments deletes vRO environments whose name starts with TestPrefix.
func (r *CleanupRunner) OrchestratorEnvironments() {
	r.applyCleanups(
		r.vROCleanupsByPrefix("vco/api/environments", "Environment", "name"),
	)
}

// OrchestratorEnvironmentRepositories deletes vRO environment repositories whose name starts with
// TestPrefix.
func (r *CleanupRunner) OrchestratorEnvironmentRepositories() {
	r.applyCleanups(
		r.vROCleanupsByPrefix(
			"vco/api/environments/repositories",
			"EnvironmentRepository",
			"name"),
	)
}

// CatalogSources deletes catalog sources whose name starts with TestPrefix.
func (r *CleanupRunner) CatalogSources() {
	r.applyCleanups(
		r.contentCleanupsByPrefix("catalog/api/admin/sources", "name"),
	)
}

// CatalogSourcesInProject deletes catalog sources in a specific project whose name starts with
// TestPrefix.
func (r *CleanupRunner) CatalogSourcesInProject(projectID string) {
	r.applyCleanups(
		r.contentCleanupsByPrefixInProject(
			"catalog/api/admin/sources", TestPrefix, projectID),
	)
}

// CustomForms fetches and deletes a custom form for the given source, if its name starts with
// TestPrefix.
func (r *CleanupRunner) CustomForms(sourceID, sourceType string) {
	fetchPath := "form-service/api/forms/fetchBySourceAndType"
	var raw CustomFormAPIModel
	resp, err := r.Client.R(fetchPath).
		SetQueryParam("formFormat", "JSON").
		SetQueryParam("formType", "requestForm").
		SetQueryParam("sourceId", sourceID).
		SetQueryParam("sourceType", sourceType).
		SetResult(&raw).
		Get(fetchPath)
	if err := r.Client.HandleAPIResponse(resp, err, []int{200, 404}); err != nil {
		r.Log.Logf("Warning: cannot fetch custom form for source %s/%s: %v", sourceID, sourceType, err)
		return
	}
	if resp.StatusCode() == 404 || raw.Id == "" {
		return
	}
	if !strings.HasPrefix(raw.Name, TestPrefix) {
		return
	}
	r.applyCleanups([]cleanupEntry{{
		label:      fmt.Sprintf("CustomForm %q (source %s)", raw.Name, sourceID),
		deletePath: "form-service/api/forms/" + raw.Id,
	}})
}

// ABXActions deletes ABX actions in the given project whose name starts with TestPrefix.
func (r *CleanupRunner) ABXActions(projectID string) {
	r.applyCleanups(
		r.contentCleanupsByPrefixInProject(
			"abx/api/resources/actions", TestPrefix, projectID),
	)
}

// CustomResourceABXActions deletes ABX actions in the given project whose name starts with
// TestCustomResourcePrefix (e.g. "Custom.AriaProviderTest.create").
func (r *CleanupRunner) CustomResourceABXActions(projectID string) {
	r.applyCleanups(
		r.contentCleanupsByPrefixInProject(
			"abx/api/resources/actions", TestCustomResourcePrefix, projectID),
	)
}

// ABXConstants deletes ABX constants (action-secrets) whose name starts with TestPrefix.
func (r *CleanupRunner) ABXConstants() {
	r.applyCleanups(
		r.contentCleanupsByPrefix("abx/api/resources/action-secrets", "name"),
	)
}

// Tags deletes iaas tags whose key starts with TestPrefix.
// When Force is true, ?ignoreUsage=true is appended to delete tags even if still in use.
func (r *CleanupRunner) Tags() {
	entries := r.contentCleanupsByPrefix("iaas/api/tags", "key")
	if r.Force {
		for i := range entries {
			entries[i].deletePath += "?ignoreUsage=true"
		}
	}
	r.applyCleanups(entries)
}

// PropertyGroups deletes property groups whose name starts with TestPrefix.
func (r *CleanupRunner) PropertyGroups() {
	r.applyCleanups(
		r.contentCleanupsByPrefix(
			"properties/api/property-groups", "name"),
	)
}

// Policies deletes policies whose name starts with TestPrefix.
func (r *CleanupRunner) Policies() {
	r.applyCleanups(
		r.contentCleanupsByPrefix("policy/api/policies", "name"),
	)
}

// Subscriptions deletes event subscriptions whose name starts with TestPrefix.
func (r *CleanupRunner) Subscriptions() {
	r.applyCleanups(
		r.contentCleanupsByPrefix("event-broker/api/subscriptions", "name"),
	)
}

// CloudTemplates deletes blueprints whose name starts with TestPrefix.
func (r *CleanupRunner) CloudTemplates() {
	r.applyCleanups(
		r.contentCleanupsByPrefix("blueprint/api/blueprints", "name"),
	)
}

// Projects deletes projects whose name starts with TestPrefix.
func (r *CleanupRunner) Projects() {
	r.applyCleanups(
		r.contentCleanupsByPrefix("project-service/api/projects", "name"),
	)
}

// CustomNamings deletes custom naming rules whose name starts with TestPrefix.
func (r *CleanupRunner) CustomNamings() {
	r.applyCleanups(
		r.contentCleanupsByPrefix("iaas/api/naming", "name"),
	)
}

// CustomResources deletes custom resource types whose displayName starts with TestPrefix.
func (r *CleanupRunner) CustomResources() {
	r.applyCleanups(
		r.contentCleanupsByPrefix(
			"form-service/api/custom/resource-types", "displayName"),
	)
}
