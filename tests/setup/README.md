# Acceptance tests setup

Terraform configuration creating the prerequisites the acceptance tests expect to already exist on
the Aria instance. Applying it writes an `env.sh` exporting every `TF_VAR_test_*` value documented
in the root [README](../../README.md#acceptance-tests).

## What it creates

| Fixture | Resource | Feeds |
|---------|----------|-------|
| 3 projects | `vra_project` | `test_project_id`, `test_project_ids` |
| ABX action | `aria_abx_action` | `test_abx_action_id`, `test_org_id` |
| Icon | `aria_icon` | `test_icon_id` |
| Cloud template, its released version and a catalog source | `vra_blueprint`, `vra_blueprint_version`, `vra_catalog_source_blueprint` | `test_catalog_item_id`, `test_catalog_item_type` |
| Secret | `restful_resource` on `platform/api/secrets` | `test_secret_id` |

Nothing is created for `test_approver_name`, the approval policy tests approve on ourselves. The
account comes from `/csp/gateway/am/api/loggedin/user` and is stripped of its domain, matching the
`USER:TOTO` form.

The secret goes through the `restful` provider because neither the aria provider (data source only)
nor the vra provider exposes a secret resource. An `aria_secret` resource would remove that provider
and the three variables it needs.

Fixtures are named `ARIA_PROVIDER_FIXTURE_*` on purpose. The `cleanup` binary sweeps everything
named `ARIA_PROVIDER_TEST*`. Fixtures therefore survive a cleanup run, which removes only the
leftovers of an interrupted test run.

## Usage

Nothing is configured here, every provider reads the environment. These exports are all that is
needed, the makefile mirrors them to the `VRA_*` and `TF_VAR_aria_*` names the vra and restful
providers expect:

```shell
export ARIA_HOST=https://some-aria-host.net
export ARIA_INSECURE=false
export ARIA_TENANT=classic # VCF 9 only, name of the VM Apps tenant
export ARIA_REFRESH_TOKEN=*****
export ARIA_ACCESS_TOKEN=*****
```

Then, from the repository root:

```shell
make testacc-setup    # terraform init && terraform apply, writes tests/setup/env.sh
make testacc-all      # the above, then the acceptance tests against those fixtures
make testacc-destroy  # tear the fixtures down
```

Terraform flags go through `ARGS`, e.g. `make testacc-setup ARGS=-auto-approve`.

To run the tests separately from a setup applied earlier:

```shell
source tests/setup/env.sh
make testacc
```

The `aria` provider is pulled from the registry. To exercise your local build instead, declare a
`dev_overrides` block in your Terraform CLI configuration.

Run `make cleanup` before `make testacc-destroy` if a test run was interrupted, a leftover
`ARIA_PROVIDER_TEST*` resource inside a fixture project blocks the project deletion. That target
loads `env.sh` on its own to scope the ABX actions and custom forms it sweeps.
