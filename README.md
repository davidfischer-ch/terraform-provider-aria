# Terraform Provider Aria (Terraform Plugin Framework)

This is the [Terraform](https://www.terraform.io) provider for VMWare's Aria Automation Platform.

The provider is [published here](https://registry.terraform.io/providers/davidfischer-ch/aria/latest).

It has been developped by the CSC Team from the IT department of the State of Geneva (Switzerland).

Please be aware that Broadcom is not responsible neither involved on this project.

_This provider is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework). See [Which SDK Should I Use?](https://developer.hashicorp.com/terraform/plugin/framework-benefits) in the Terraform documentation for additional information._


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25


## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```


## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.


## Using the provider

Fill this in for each provider


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your
machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary
in the `$GOPATH/bin` directory.

To generate or update documentation, run `make docs`.
To format the code run `make fmt`.
To do both in one step, run `make tidy`.

### Pre-push checks

`make check` runs the same steps as the CI **Tests** workflow (build, vet, lint, generate, unit
tests) against your working tree. Run it before pushing to catch failures locally. The unit tests
step provides its own dummy `ARIA_HOST` and `ARIA_REFRESH_TOKEN`, no setup required, and writes an
HTML coverage report to `bin/coverage.html` (see [Unit tests](#unit-tests) below).

```shell
make check                    # all steps
make check ARGS="build lint"  # only the named steps
```

The script lives at `scripts/check.sh` and can be run directly. It installs the pinned golangci-lint
release into `./bin` when absent, matching the version CI uses.

### Linting

Requires golangci-lint v2 (the `.golangci.yml` config uses the v2 schema).

```shell
make lint
```

### Unit tests

```shell
make test
```

`make test` sets dummy `ARIA_HOST`/`ARIA_REFRESH_TOKEN` values for you (no real API is called), and
writes an HTML coverage report to `bin/coverage.html` (open it in a browser). To run tests with your
own values or flags, set the environment variables yourself and run `go test ./...` directly.

### Acceptance tests

Acceptance tests create and destroy real resources on a live Aria instance.

These exports are all the makefile targets need, they derive the rest:

```shell
export ARIA_HOST=https://some-aria-host.net
export ARIA_INSECURE=false
export ARIA_TENANT=classic # VCF 9 only, name of the VM Apps tenant, uses the VCF 9 API token flow
export ARIA_REFRESH_TOKEN=*****
export ARIA_ACCESS_TOKEN=***** # If you have one, not required
```

The tests also expect a set of resources to already exist on the instance. The
[tests/setup](tests/setup) Terraform configuration creates and manages them, writing the matching
`TF_VAR_test_*` exports to `tests/setup/env.sh`, see its [README](tests/setup/README.md):

```shell
make testacc-setup    # create the prerequisites, writes tests/setup/env.sh
make testacc-all      # the above, then the acceptance tests against those prerequisites
make testacc-destroy  # tear them down
```

That configuration uses the released provider. `DEV=1` builds the one of this worktree and points
Terraform at it instead, the only way to exercise unreleased changes:

```shell
make testacc-setup DEV=1
```

Then, once you have them available, run:

```shell
make testacc
```

`make testacc` loads `tests/setup/env.sh` when it exists. Export the `TF_VAR_test_*` values by hand
only if you manage the prerequisites yourself.

`make testacc` runs the unit tests first (`make test`, no live API needed), then the acceptance
tests against your Aria instance. This fails fast on a broken unit test before spending time on the
slower acceptance run.

`TESTACC_RUN` narrows the acceptance run to a regexp matched against test names. Go selects
functions, not files, but one file's tests usually share a prefix:

```shell
make testacc TESTACC_RUN=TestAccIconDataSource  # a single test
make testacc TESTACC_RUN='TestAccIcon.*'        # every test of icon_resource_acc_test.go
```

The unit run that gates it still runs in full. `TEST_RUN` narrows that one, on `make test` or
`make check`:

```shell
make test TEST_RUN=TestCustomResourceModelToAPI
```

Both default to the whole suite, and coverage is partial whenever a run is narrowed.

The `TF_VAR_test_catalog_item_*` variables point to an existing catalog item whose icon and custom
form **will be modified** by the tests.

### Cleaning up test resources

If an acceptance test run is interrupted or fails mid-way, orphaned resources may remain on the
Aria instance. The `cleanup` binary sweeps all resources whose names follow the `ARIA_PROVIDER_TEST`
prefix convention used by the test suite.

`make cleanup` builds it, loads `tests/setup/env.sh` when present, and runs it. Flags go through
`ARGS`:

```shell
make cleanup ARGS=-help
make cleanup ARGS=-dry-run  # preview what would be deleted without touching the API
make cleanup                # delete everything
make cleanup ARGS=-force    # also bypass vRO dependency checks and tag usage locks
```

Build and run it by hand instead if you prefer:

```shell
go build -o bin/cleanup ./cmd/cleanup/
bin/cleanup -dry-run
```

The `TF_VAR_test_project_id`, `TF_VAR_test_catalog_item_id`, and `TF_VAR_test_catalog_item_type`
environment variables are reused from the acceptance test setup above to scope ABX actions
and custom forms cleanup.

Run it before `make testacc-destroy`, a leftover `ARIA_PROVIDER_TEST*` resource inside a
fixture project blocks the project deletion.
