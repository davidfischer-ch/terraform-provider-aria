# Changelog

## Release v0.6.10 (2025-02-21)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.9...v0.6.10

### Features

* Resource `aria_custom_form`: Seamlessly use (and overwrite) existing form

### Fix and enhancements

* Merge dependabot dependencies update requests

## Release v0.6.9 (2025-02-21)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.8...v0.6.9

### Features

* Resource `aria_orchestrator_action`: Serialize write requests to prevent vRO deadlocks

### Fix and enhancements

* Merge dependabot dependencies update requests

## Release v0.6.8 (2025-02-17)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.7...v0.6.8

### Features

* Resource `aria_custom_naming`: Add `counters` nested attribute (fix #38 counters are reset on UPDATE)
* Provider config: Add `ok_api_calls_log_level` (defaulting to `TRACE`)
+ Provider config: Add `ko_api_calls_log_level` (defaulting to `DEBUG`)

### Fix and enhancements

* SDK: Differentiate from/to API content and implement read after create and/or update if necessary
* SDK: API client's tries to indent response Body (if JSON) to make it more readable in logs
* SDK: API client's do not set apiVersion header if version is ""
* SDK: Refactor API client (add more functions) to reduce code boilerplate (make code DRY)
* Docs: Drop caution since #114 (random catalog item's download error...) is fixed
* CI: Fix linters deprecations

## Release v0.6.7 (2025-02-11)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.6...v0.6.7

### Features

* Resource `aria_catalog_source`: Implement Cloud Templates & ABX Actions use-cases

## Release v0.6.6 (2025-02-10)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.5...v0.6.6

### Features

* Add resource `aria_orchestrator_environment`
* Add resource `aria_orchestrator_environment_repository`

## Release v0.6.5 (2025-02-06)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.4...v0.6.5

### Features

* Add resource `aria_orchestrator_task`
* Resource `aria_catalog_source`, if `wait_imported = true` then the polling function will (up to 15 minutes) :
  * Proactively wait for catalog items to be available (by doing the equivalent of save & import multiple times)
  * Immediately return an error if there are any *hard* errors (not related to catalog item's availability)

## Release v0.6.4 (2025-01-31)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.3...v0.6.4

### Features

* Add data source `aria_orchestrator_configuration`
* Add resource `aria_tag`

### Fix and enhancements

* Mark `aria_policy`'s `scope_criteria` as immutable
* Fix `aria_icon`'s `keep_on_destroy` attribute cannot be updated
* SDK: Simplify code (do not return & manage empty diagnostics)
* Doc: Improve it with more examples & declare immutable attributes

## Release v0.6.3 (2025-01-30)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.2...v0.6.3

### Features

* Add data source `aria_integration`
* Add resource `aria_catalog_source` (with an open isse, see [#114](https://github.com/davidfischer-ch/terraform-provider-aria/issues/114)
* Add resource `aria_orchestrator_configuration`
* Add resource `aria_policy`

### Fix and enhancements

* Validate `aria_orchestrator_category`'s `type` attribute
* SDK: Make `AttributeTypes` a func of `*Model` struct (as per in doc)
* Merge dependabot dependencies update requests

## Release v0.6.2 (2024-12-11)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.1...v0.6.2

### Minor compatibility breaks

* SDK: API client's `ReadIt` method also return API call's response (e.g. to let you read headers)
* SDK: Drop `ActionParameter<API>Model` (replaced by `Parameter<API>Model`)

### Features

* Add resource `aria_orchestrator_workflow`

### Fix and enhancements

* API client's `ReadIt` method: Add `readPath` optional argument
* Deduplicate parameter models and schemas
* Change icon used by tests to prevent a conflict with icons we use at OCSIN
* Fix Changelog

## Release v0.6.1 (2024-12-05)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.6.0...v0.6.1

### Features

The platform's icon API endpoint is not thread safe neither prevent the deletion of icons used by the catalog items ...

So :

* Implement `aria_catalog_item_icon`'s soft deletion (opt-in with `keep_on_destroy = true`)
* Prevent Aria internal errors when manipulating multiple icons by using mutexes on `aria_icon`'s resource CRUD functions

### Fix and enhancements

* Cover `aria_catalog_item_icon` with tests

## Release v0.6.0 (2024-12-04)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.5.6...v0.6.0

### Major compatibility breaks

* Icons resources are now declared with `path` & (optional) `hash` instead of `content`

### Features

* Add resource `aria_catalog_item_icon`
* Make `aria_icon` compatible with any kind of image

### Fix and enhancements

* Don't store `aria_icon`'s content inside the state (only path & its SHA-256 checksum)

## Release v0.5.6 (2024-12-03)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.5.5...v0.5.6

### Features

* Implement `aria_orchestrator_action`'s forced deletion (opt-in with `force_delete = true`)
* All resources : Retry deletion upon conflicts error (409's) to try to converge (up to 5 times)

#### Deletion

When deleting a resource, the conflict error (409) is potentially an error that will be solved by the deletion of other resources made in parallel in the same apply.

Ideally, one should declare the dependencies with the `depends_on` or by using attribute's of a resource to configure another. However sometimes its not practical or even not possible (e.g. creating resources with a `for_each` loop). Hence Terraform will execute all delete operations in parallel.

This is why the delete function now retries the delete operation up to 5 times to try to converge to desired state. This is not optimal but at least working for the common use cases.

If this "magic number" of 5 (or the delay) has to be tuned per resource, please open an issue.

### Fix and enhancements

* Document import of resources when available

## Release v0.5.5 (2024-11-29)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.5.4...v0.5.5

### Features

* Add resource `aria_orchestrator_action`
* Add resource `aria_orchestrator_category`

### Fix and enhancements

* Merge dependabot merge requests:
  * Bump github.com/go-resty/resty/v2 from 2.15.3 to 2.16.2
  * Bump github.com/hashicorp/terraform-plugin-docs from 0.19.4 to 0.20.1
  * Bump github.com/hashicorp/terraform-plugin-framework-validators from 0.13.0 to 0.15.0
  * Bump github.com/hashicorp/terraform-plugin-testing from 1.10.0 to 1.11.0
* Make code consistent

## Release v0.5.4 (2024-10-08)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.5.3...v0.5.4

### Features

* Add ressource `aria_custom_form`

## Release v0.5.3 (2024-09-30)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.5.2...v0.5.3

### Fix and enhancements

* Fix `aria_custom_resource`'s `resource_type` attribute change must force recreation of the resource #66
  * Thanks to @selknsi
* Bump `github.com/go-resty/resty/v2` from `2.14.0` to `2.15.1` #70
* Drop "example" function

## Release v0.5.2 (2024-09-02)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.5.1...v0.5.2

### Fix and enhancements

* Fix `aria_subscription`'s project scoping: Was scoped to 0 projects instead of being unscoped
* Fix `aria_sbuscription`'s `owner_id`: It may change, cannot use previous known value

## Release v0.5.1 (2024-08-26)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.5.0...v0.5.1

### Fix and enhancements

* Fix UPDATE on `aria_custom_resources` when projectID is empty (by omitting this field on update, crazy POST != full object)

## Release v0.5.0 (2024-08-22)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.4.1...v0.5.0

This is not a major release because its still a v0.x and incrementing to v1 will be done when used
in production for at least few months.

### Major compatibility breaks

* Properties default value: Replace custom state storage encoding format by JSON (see migrations)

### Migrations

You have to manipulate the state and replace default values:

* String: `some text` by `\"some text\"`
* Float: `%!s(float64=10240)` by `10240`
* ... (more jsonencoded values here) ...

Here is the error:

```
│ Error: Invalid JSON String Value
│
│   with aria_custom_resource.redis,
│   on redis.tf line 80, in resource "aria_custom_resource" "redis":
│   80: resource "aria_custom_resource" "redis" {
│
│ A string value was provided that is not valid JSON string format (RFC 7159).
│
│ Given Value: %!s(float64=10240)
╵
```

### Features

* Resource `aria_resource_action`: Make it compatible with custom resources, fix #19

#### Work in progress

* Add work in progress resource `aria_project`
* Add work in progress resource `aria_cloud_template_v1`

### Fix and enhancements

* Update `aria_custom_resource` example (doc)
* Merge dependabot requests (terraform plugin testing, resty, ...)
* Fix #61 by replacing propertie's default value custom state storage encoding format with JSON

### Library

#### Features

* Add `AriaClient` in replacement of `AriaClientConfig`
  * Add `Mutex` attribute, used internally by the resources
  * Add `ReadIt` method to deduplicate resources Read methods
  * Add a bunch of other methods
* Library: Add `RWMutexKV` (key/value store for arbitrary read-write mutexes)
* Add `Model` interface with utility methods (declaration) to be able to factorize code

#### Fix and enhancements

* Implement utility methods on models that are exposed as resources
* Refactor resources to use `AriaClient` capabilities

## Release v0.4.1 (2024-07-25)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.4.0...v0.4.1

### Fix and enhancements

* Resource `aria_resource_action`: Attribute `criteria` is now normalized with `jsontypes.Normalized`
* Resource `aria_resource_action`: Attribute `form_definition.form` is now normalized with `jsontypes.Normalized`

See https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework-jsontypes@v0.1.0/jsontypes#Normalized:

Semantic equality logic is defined for Normalized such that inconsequential differences between JSON strings are ignored (whitespace, property order, etc). If you need strict, byte-for-byte, string equality, consider using ExactType.

## Release v0.4.0 (2024-07-23)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.3.1...v0.4.0

### Minor compatibility breaks

* Resource `aria_custom_resource`: Attribute `properties` is now a map of name to property
* Resource `aria_property_group`: Attribute `properties` is now a map of name to property

### Features

* Diff of `aria_custom_resource` and `aria_property_group` should be more readable (+/- & ~)
* Managing `aria_custom_resource` and `aria_property_group` will be immune to Terraform "diff" after apply errors

## Release v0.3.1 (2024-07-23)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.3.0...v0.3.1

### Features

* Resource `aria_resource_action`: Add attribute `criteria`
* Resource `aria_resource_action`: Add attribute `form_definition`

## Release v0.3.0 (2024-07-16)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.7...v0.3.0

### Minor compatibility breaks

* Make some parameters of `aria_custom_resource.properties` mandatory (`encrypted`, `read_only`, `recreated_on_update`)
* Make `aria_custom_resource.properties` unordered (from ordered, see #37)

### Features

* Add resource `aria_property_group`
* Make API request info available in `TRACE` (when successful) or `DEBUG` (when failed)

### Fix and enhancements

* Make code DRY
* Test imports in acceptance tests
* Add TODOs for the future

## Release v0.2.7 (2024-07-11)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.6...v0.2.7

### Fix and enhancements

* Resource `aria_custom_naming`: Manage changing templates (update in place)
* Add CAUTION section in `aria_custom_naming` resource description
* Log API call details only in case of error

## Release v0.2.6 (2024-07-11)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.5...v0.2.6

### Fix and enhancements

* Resource `aria_resource_action`: Omit `project_id` when empty (JSON marshaling)

## Release v0.2.5 (2024-07-10)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.4...v0.2.5

### Features

* Add resource `aria_custom_naming`

### Fix and enhancements

* Include request in debug info + make it multiline instead of multi log lines
* Define Aria APIs versions and explicitly define `apiVersion` when making API requests
* Omit empty `id` attribute when making API requests
* Use `int32` when `int64` is overkill to save some cache
* Upgrade go modules
* Instruct Terraform to replace a `aria_resource_action` when its `project_id` is changed

## Release v0.2.4 (2024-07-09)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.3...v0.2.4

### Fix and enhancements

* Make `aria_resource_action.runnable_item.name` required

## Release v0.2.3 (2024-07-09)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.2...v0.2.3

### Features

* Add resource `aria_custom_resource` with some missing features
* Add resource `aria_resource_action` limited to natives types (for now)

## Release v0.2.2 (2024-06-24)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.1...v0.2.2

### Features

* Add data source `aria_secret`

## Release v0.2.1 (2204-06-24)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.0...v0.2.1

### Features

* Resource `aria_abx_action`: Add attributes `cpu_shares`, `deployment_timeout_seconds`, `shared`, `system`, `async_deployed`

## Release v0.2.0 (2024-06-21)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.1.2...v0.2.0

### Minor compatibility breaks

* Rename resource `aria_abx_secret` to `aria_abx_sensitive_constant`

## Release v0.1.2 (2024-06-19)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.1.1...v0.1.2

### Fix and enhancements

* Resource `aria_abx_action`: Fix conversion from/to "" <-> "auto" (#9)
* Documentation: Update example code and description of resources

## Release v0.1.1 (2024-06-19)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.1.0...v0.1.1

### Features

* Add resource `aria_abx_action`
* Add resource `aria_subscription`
* Add data source `aria_catalog_type`

### Fix and enhancements

* Documentation: Update section about acceptance testing
* Dependencies: Upgrade terraform-plugin-framework to 1.9.0
* Code Style: Favor self over one letter names

## Release v0.1.0 (2024-06-05)

Initial release.
