# Version information from git
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
CLEAN_VERSION = $(patsubst v%,%,$(VERSION))

# Version extraction and manipulation
CURRENT_VERSION := $(shell v=$$(git tag --list "v[0-9]*.[0-9]*.[0-9]*" --sort=-v:refname | head -n1 | sed 's/^v//' | tr -d '\n'); if [ -z "$$v" ]; then echo "0.0.0"; else echo "$$v"; fi)
MAJOR := $(shell echo $(CURRENT_VERSION) | cut -d. -f1)
MINOR := $(shell echo $(CURRENT_VERSION) | cut -d. -f2)
PATCH := $(shell echo $(CURRENT_VERSION) | cut -d. -f3)
NEW_PATCH := $(shell echo $$(($(PATCH) + 1)))
NEW_MINOR := $(shell echo $$(($(MINOR) + 1)))
NEW_MAJOR := $(shell echo $$(($(MAJOR) + 1)))
NEXT_PATCH_VERSION := $(MAJOR).$(MINOR).$(NEW_PATCH)
NEXT_MINOR_VERSION := $(MAJOR).$(NEW_MINOR).0
NEXT_MAJOR_VERSION := $(NEW_MAJOR).0.0
# Allow pre-release and build metadata in version validation
VERSION_VALID := $(shell echo $(CURRENT_VERSION) | grep -E '^[0-9]+\.[0-9]+\.[0-9]+(-[A-Za-z0-9\.-]+)?(\+[A-Za-z0-9\.-]+)?$$' >/dev/null && echo "true" || echo "false")
IS_DIRTY := $(shell git diff-index --quiet HEAD -- || echo "true")

# Default target
.DEFAULT_GOAL := help

.PHONY: help
# Help message
help:
	@echo "Go Patterns:"
	@echo "  test                 - Run all tests"
	@echo "  benchmark            - Run all benchmarks"
	@echo "  test-one  [pkg=]      - Run tests for specific package (make test-one pkg=optargs)"
	@echo "  benchmark-one [pkg=]  - Run benchmarks for specific package (make benchmark-one pkg=optargs)"
	@echo ""
	@echo "  Maintenance Targets:"
	@echo "  help                 - Show this help message"
	@echo ""
	@echo "  Version Management:"
	@echo "  version-bump-patch    - Bump patch version (1.2.3 -> 1.2.4)"
	@echo "  version-bump-minor    - Bump minor version (1.2.3 -> 1.3.0)"
	@echo "  version-bump-major    - Bump major version (1.2.3 -> 2.0.0)"
	@echo "  version-set           - Set specific version (make version-set VERSION=1.2.3)"
	@echo "  push          		   - Push version tags to remote"




# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./... 

.PHONY: benchmark
benchmark:
	@echo "Running benchmarks..."
	go test ./... -bench=. -benchmem

.PHONY: test-one
test-one:
	@echo "Running tests for ${pkg}..."
	go test -v  ./pkg/${pkg}/... 

.PHONY: benchmark-one
benchmark-one:
	@echo "Running benchmarks for ${pkg}..."
	go test ./pkg/${pkg}/... -bench=. -benchmem


# Validate version format
.PHONY: validate-version
validate-version:
	@echo "[DEBUG] CURRENT_VERSION: '$(CURRENT_VERSION)'"
	@if [ "$(VERSION_VALID)" != "true" ]; then \
		echo "Error: Current version '$(CURRENT_VERSION)' does not follow semantic versioning format (major.minor.patch)"; \
		exit 1; \
	fi


# Bump patch version
.PHONY: version-bump-patch
version-bump-patch: validate-version
	@if [ "$(IS_DIRTY)" = "true" ]; then \
		echo "Warning: You have uncommitted changes."; \
		read -p "Continue anyway? [y/N] " confirm; \
		if [ "$$confirm" != "y" ]; then exit 1; fi; \
	fi
	@echo "Bumping patch version: $(CURRENT_VERSION) -> $(NEXT_PATCH_VERSION)"
	@git tag -a v$(NEXT_PATCH_VERSION) -m "Bump patch version to $(NEXT_PATCH_VERSION)"
	@echo "Tagged with v$(NEXT_PATCH_VERSION)"

# Bump minor version
.PHONY: version-bump-minor
version-bump-minor: validate-version
	@if [ "$(IS_DIRTY)" = "true" ]; then \
		echo "Warning: You have uncommitted changes."; \
		read -p "Continue anyway? [y/N] " confirm; \
		if [ "$$confirm" != "y" ]; then exit 1; fi; \
	fi
	@echo "Bumping minor version: $(CURRENT_VERSION) -> $(NEXT_MINOR_VERSION)"
	@git tag -a v$(NEXT_MINOR_VERSION) -m "Bump minor version to $(NEXT_MINOR_VERSION)"
	@echo "Tagged with v$(NEXT_MINOR_VERSION)"

# Bump major version
.PHONY: version-bump-major
version-bump-major: validate-version
	@if [ "$(IS_DIRTY)" = "true" ]; then \
		echo "Warning: You have uncommitted changes."; \
		read -p "Continue anyway? [y/N] " confirm; \
		if [ "$$confirm" != "y" ]; then exit 1; fi; \
	fi
	@echo "Bumping major version: $(CURRENT_VERSION) -> $(NEXT_MAJOR_VERSION)"
	@git tag -a v$(NEXT_MAJOR_VERSION) -m "Bump major version to $(NEXT_MAJOR_VERSION)"
	@echo "Tagged with v$(NEXT_MAJOR_VERSION)"

# Set specific version
.PHONY: version-set
version-set:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION not specified. Usage: make version-set VERSION=1.2.3"; \
		exit 1; \
	fi
	@if [ "$(IS_DIRTY)" = "true" ]; then \
		echo "Warning: You have uncommitted changes."; \
		read -p "Continue anyway? [y/N] " confirm; \
		if [ "$$confirm" != "y" ]; then exit 1; fi; \
	fi
	@echo "Setting version to $(VERSION)"
	@git tag -a v$(VERSION) -m "Set version to $(VERSION)"
	@echo "Tagged with v$(VERSION)"

# Push version tags to remote
.PHONY: push
push:
	@git push origin --tags

