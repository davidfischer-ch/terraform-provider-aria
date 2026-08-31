default: testacc

# Run the local pre-push checks that mirror the CI "Tests" workflow
# (build, vet, lint, generate, unit tests). Pass steps as args, e.g.
# `make check ARGS="build lint"`.
.PHONY: check
check:
	./scripts/check.sh $(ARGS)

# Format Go source (-s also simplifies code, e.g. dropping redundant type conversions).
.PHONY: fmt
fmt:
	gofmt -l -s -w .

# Regenerate docs/ from the schema, and format examples/ (see main.go's go:generate directives).
.PHONY: docs
docs:
	go generate ./...

# Format code and regenerate docs in one go.
.PHONY: tidy
tidy: fmt docs

# Lint (installs the pinned golangci-lint release into ./bin when absent, matching CI).
.PHONY: lint
lint:
	./scripts/check.sh lint

# Narrow a run to the test names the regexp matches. Go has no per-file selection, tests are
# functions, but one file's tests usually share a prefix:
#   make test TEST_RUN=TestCustomResourceModelToAPI
#   make testacc TESTACC_RUN=TestAccIconDataSource
#   make testacc TESTACC_RUN='TestAccIcon.*'   # every test of icon_resource_acc_test.go
# They are separate variables on purpose: narrowing the acceptance run must not silently empty the
# unit run that gates it.
TEST_RUN ?=
TESTACC_RUN ?= ^TestAcc

# Run unit tests (dummy ARIA_HOST/ARIA_REFRESH_TOKEN are set for you, no real API is called).
.PHONY: test
test:
	TEST_RUN='$(TEST_RUN)' ./scripts/check.sh test

# Run unit tests first (fast, no live API needed): a broken unit test then fails before any time
# is spent on the slow, real-API acceptance run. Go interleaves *_unit_test.go and *_acc_test.go
# files alphabetically within a package; without this split a broken unit test can sit behind
# several acceptance tests instead of failing immediately.
#
# The TF_VAR_test_* values written by testacc-setup are loaded when present, running the tests
# against those fixtures without sourcing anything by hand.
.PHONY: testacc
testacc: test
	if [ -f ./$(TESTACC_SETUP_DIR)/env.sh ]; then . ./$(TESTACC_SETUP_DIR)/env.sh; fi; \
	TF_ACC=1 go test ./... -v -run '$(TESTACC_RUN)' $(TESTARGS) -timeout 120m

# The setup configuration mirrors ARIA_* to what the vra and restful providers expect, the aria
# provider reads ARIA_* on its own. See tests/setup/README.md. These are exported rather than
# prefixed to the recipes, which would print the tokens and expose them in the process list.
TESTACC_SETUP_DIR = tests/setup
ARIA_INSECURE ?= false
export VRA_URL = $(ARIA_HOST)
export VRA_REFRESH_TOKEN = $(ARIA_REFRESH_TOKEN)
export VRA_INSECURE = $(ARIA_INSECURE)
export TF_VAR_aria_host = $(ARIA_HOST)
export TF_VAR_aria_access_token = $(ARIA_ACCESS_TOKEN)
export TF_VAR_aria_insecure = $(ARIA_INSECURE)

# Create the prerequisites the acceptance tests expect and write tests/setup/env.sh.
# Pass Terraform flags as args, e.g. `make testacc-setup ARGS=-auto-approve`.
.PHONY: testacc-setup
testacc-setup:
	cd $(TESTACC_SETUP_DIR) && terraform init
	cd $(TESTACC_SETUP_DIR) && terraform apply $(ARGS)

# Destroy the prerequisites. Run bin/cleanup first if an acceptance run was interrupted, a leftover
# ARIA_PROVIDER_TEST* resource inside a fixture project blocks the project deletion.
.PHONY: testacc-destroy
testacc-destroy:
	cd $(TESTACC_SETUP_DIR) && terraform destroy $(ARGS)

# Create the prerequisites, then run the acceptance tests against them.
.PHONY: testacc-all
testacc-all: testacc-setup
	$(MAKE) testacc

# Sweep the ARIA_PROVIDER_TEST* resources left over by an interrupted acceptance run. The
# TF_VAR_test_* values written by testacc-setup scope the ABX actions and custom forms it looks at,
# they are loaded when present. Pass flags as args, e.g. `make cleanup ARGS=-dry-run`.
.PHONY: cleanup
cleanup:
	go build -o bin/cleanup ./cmd/cleanup/
	if [ -f ./$(TESTACC_SETUP_DIR)/env.sh ]; then . ./$(TESTACC_SETUP_DIR)/env.sh; fi; \
	bin/cleanup $(ARGS)
