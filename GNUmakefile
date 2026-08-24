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

# Run unit tests (dummy ARIA_HOST/ARIA_REFRESH_TOKEN are set for you, no real API is called).
.PHONY: test
test:
	./scripts/check.sh test

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
