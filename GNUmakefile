default: testacc

# Run the local pre-push checks that mirror the CI "Tests" workflow
# (build, vet, lint, generate, unit tests). Pass steps as args, e.g.
# `make check ARGS="build lint"`.
.PHONY: check
check:
	./scripts/check.sh $(ARGS)

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
