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
| Orchestrator environment | `aria_orchestrator_environment` | `test_environment_id`, `test_environment_name`, `test_runtime` |
| Orchestrator category and workflow | `aria_orchestrator_category`, `aria_orchestrator_workflow` | `test_workflow_id`, `test_workflow_name` |

Nothing is created for `test_approver_name`, the approval policy tests approve on ourselves. The
account comes from `/csp/gateway/am/api/loggedin/user` and is stripped of its domain, matching the
`USER:TOTO` form.

The secret goes through the `restful` provider because neither the aria provider (data source only)
nor the vra provider exposes a secret resource. An `aria_secret` resource would remove that provider
and the three variables it needs.

The environment is built on the runtime of the platform, `python:3.11` on VCF 9 and `python:3.10`
on Aria Automation 8.x. The runtimes an appliance offers differ from one version to another, and an
unavailable one is reported as `MISSING_RUNTIME` in the validation message of every resource using
it. Override the default when neither is one of yours:

```shell
export TF_VAR_test_runtime=python:3.12
```

The workflow waits for its service broker import, up to fifteen minutes. A resource action backed
by a workflow needs the broker to know it, and that import is asynchronous: paying for it once here
is what lets the tests consuming this workflow start from an imported one. An appliance whose
import never completes therefore fails the setup, which is the intended report, the alternative
being every run discovering it separately.

Fixtures are named `ARIA_PROVIDER_FIXTURE_*` on purpose. The `cleanup` binary sweeps everything
named `ARIA_PROVIDER_TEST*`. Fixtures therefore survive a cleanup run, which removes only the
leftovers of an interrupted test run.

## Usage

The aria provider reads the environment on its own. The vra and restful providers are configured
explicitly, from variables the makefile mirrors under `TF_VAR_aria_*`. These exports are all that is
needed:

```shell
export ARIA_HOST=https://some-aria-host.net
export ARIA_INSECURE=false
export ARIA_TENANT=classic # VCF 9 only, name of the VM Apps tenant
export ARIA_REFRESH_TOKEN=*****
```

Then, from the repository root:

```shell
make testacc-setup    # terraform init && terraform apply, writes tests/setup/env.sh
make testacc-all      # the above, then the acceptance tests against those fixtures
make testacc-destroy  # tear the fixtures down
```

Terraform flags go through `ARGS`, e.g. `make testacc-setup ARGS=-auto-approve`.

To run the tests separately from a setup applied earlier, `make testacc` picks `env.sh` up on its
own:

```shell
make testacc
```

The `aria` provider is pulled from the registry. `DEV=1` builds the one of this worktree and points
Terraform at it through a `dev_overrides` block, the only way to exercise unreleased changes:

```shell
make testacc-setup DEV=1
```

Terraform then prints a "Provider development overrides are in effect" warning naming the binary.
The generated CLI configuration lives in `bin/`, nothing outside the repository is touched and the
override lasts for that command alone. The lock file keeps pinning the released version. Stay on
`DEV=1` for `make testacc-destroy` too: state written by the worktree build would otherwise be read
back by the registry one.

The acceptance tests need none of this, they load the provider in-process.

Run `make cleanup` before `make testacc-destroy` if a test run was interrupted, a leftover
`ARIA_PROVIDER_TEST*` resource inside a fixture project blocks the project deletion. That target
loads `env.sh` on its own to scope the ABX actions and custom forms it sweeps.
