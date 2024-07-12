# Changelog

## Release v0.2.7 (2024-07-11)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.7...v0.2.6

### Fix and enhancements

* Resource `aria_custom_naming`: Manage changing templates (update in place)
* Add CAUTION section in `aria_custom_naming` resource description
* Log API call details only in case of error

## Release v0.2.6 (2024-07-11)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.6...v0.2.5

### Fix and enhancements

* Resource `aria_resource_action`: Omit `project_id` when empty (JSON marshaling)

## Release v0.2.5 (2024-07-10)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.5...v0.2.4

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

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.4...v0.2.3

### Fix and enhancements

* Make `aria_resource_action.runnable_item.name` required

## Release v0.2.3 (2024-07-09)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.3...v0.2.2

### Features

* Add resource `aria_custom_resource` with some missing features
* Add resource `aria_resource_action` limited to natives types (for now)

## Release v0.2.2 (2024-06-24)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.2...v0.2.1

### Features

* Add data source `aria_secret`

## Release v0.2.1 (2204-06-24)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.1...v0.2.0

### Features

* Resource `aria_abx_action`: Add attributes `cpu_shares`, `deployment_timeout_seconds`, `shared`, `system`, `async_deployed`

## Release v0.2.0 (2024-06-21)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.2.0...v0.1.2

### Minor compatibility breaks

* Rename resource `aria_abx_secret` to `aria_abx_sensitive_constant`

## Release v0.1.2 (2024-06-19)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.1.2...v0.1.1

### Fix and enhancements

* Resource `aria_abx_action`: Fix conversion from/to "" <-> "auto" (#9)
* Documentation: Update example code and description of resources

## Release v0.1.1 (2024-06-19)

Diff: https://github.com/davidfischer-ch/terraform-provider-aria/compare/v0.1.1...v0.1.0

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
