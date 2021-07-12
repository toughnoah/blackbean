.PHONY: all
all: init fmt vet test install clean

define go-init
@export GOPROXY=https://goproxy.io
@go mod tidy
endef

.PHONY: init
init: ## go mod tidy
	$(call go-init)

fmt: ## Run go fmt against code.
	@go fmt ./...

vet: ## Run go vet against code.
	@go vet ./...

test: init ## Run go test against code.
	$(call go-init)
	@mv ${HOME}/.blackbean.yaml ${HOME}/blackbean.yaml
	@go test ./... -v --coverprofile=cover.out
	@go tool cover -func=cover.out
	@mv ${HOME}/blackbean.yaml ${HOME}/.blackbean.yaml

install: ## build blackbean.
	$(call go-init)
	@go build -o ${HOME}/bin/blackbean

clean: ## Remove test report.
	@rm -f *.out
