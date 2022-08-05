
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

APP_NAME=kubercert

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

GO ?= go
GO_MD2MAN ?= go-md2man

all: manager

major:  ## Set major version
	@#git tag -a v0.0.1 -m 'version v0.0.1' && git push --tags
	@git tag $$(svu major)
	@git push --tags
	@#goreleaser --rm-dist

minor:  ##  Set minor version
	@git tag $$(svu minor)
	@git push --tags
	@#goreleaser --rm-dist

patch:  ##  Set patch version
	@git tag $$(svu patch)
	@git push --tags
	@#goreleaser --rm-dist

# Run tests
test: generate fmt vet manifests
	go test ./... -coverprofile cover.out

dep: ## Get the dependencies
	@$(GO) get -v -d ./...

update: ## Get and update the dependencies
	@$(GO) get -v -d -u ./...

tidy: ## Clean up dependencies
	@$(GO) mod tidy

vendor: dep ## Create vendor directory
	@$(GO) mod vendor

build: # Build k3ctl binary
	@#go build -o ${APP_NAME} cli/main.go
	@#go build -o ${APP_NAME} main.go
	@#go build -o ${APP_NAME} -ldflags="-X main.version=0.0.2" main.go
	@goreleaser build --snapshot --single-target --rm-dist -o .

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...
