# Changelog

## Release v0.4.1 (2024-07-25)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.4.1...v0.4.0

### Fix and enhancements

* Resource `aria_resource_action`: Attribute `criteria` is now normalized with `jsontypes.Normalized`
* Resource `aria_resource_action`: Attribute `form_definition.form` is now normalized with `jsontypes.Normalized`

See https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework-jsontypes@v0.1.0/jsontypes#Normalized:

Semantic equality logic is defined for Normalized such that inconsequential differences between JSON strings are ignored (whitespace, property order, etc). If you need strict, byte-for-byte, string equality, consider using ExactType.

## Release v0.4.0 (2024-07-23)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.4.0...v0.3.1

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
