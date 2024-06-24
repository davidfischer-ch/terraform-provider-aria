# Changelog

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
