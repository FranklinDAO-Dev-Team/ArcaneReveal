#################
# golangci-lint #
#################
golangci_version=v1.56.2

.PHONY: lint-install lint lint-fix setup setup-vendor setup-vendor-docker

# Setup
setup:
	@echo "--> Setting up seismic circuits"
	@cd seismic/circuit && yarn && yarn ptau:up && yarn build
	@echo "--> Setting up seismic circuits done"
	@echo "--> Installing modvendor"
	@go install github.com/goware/modvendor@latest
	@echo "--> Vendor cgo deps"
	@go mod vendor && modvendor -copy="**/*.c **/*.h **/*.proto **/*.a **/dummy.go **/*.dylib **/*.so" -v
	@echo "--> Finished vendoring cgo deps"

setup-vendor-docker:
	@go install github.com/goware/modvendor@latest
	@GO111MODULE=on go mod vendor
	@echo "--> Vendor cgo deps"
	@go mod vendor && modvendor -copy="**/*.c **/*.h **/*.proto **/*.a **/dummy.go **/*.dylib **/*.so" -v
	@echo "--> Finished vendoring cgo deps"

# Install golangci-lint if not present
lint-install:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "--> golangci-lint is not installed. Installing..."; \
		@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version); \
	fi

# Run golangci-lint on the project
lint:
	@$(MAKE) lint-install
	@echo "--> Running golangci-lint"
	@golangci-lint run  --timeout=10m --concurrency 8 -v


# Run golangci-lint on the project and fix any issues
lint-fix:
	@$(MAKE) lint-install
	@echo "--> Running golangci-lint"
	@golangci-lint run  --timeout=10m --fix --concurrency 8 -v


