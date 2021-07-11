.PHONY: all
all: init fmt vet test install clean

define go-init
@if [ ! -f "${HOME}/.blackbean.yaml" ]; \
then \
touch ${HOME}/.blackbean.yaml; \
else \
echo "go init"; \
fi
go mod tidy
endef

.PHONY: init
init: ## Create blackbean config file and run go mod tidy
	$(call go-init)

fmt: ## Run go fmt against code.
	@go fmt ./...

vet: ## Run go vet against code.
	@go vet ./...

test: ## Run go test against code.
	@mv ${HOME}/.blackbean.yaml ${HOME}/blackbean.yaml
	@go test ./... -v --coverprofile=cover.out
	@go tool cover -func=cover.out
	@mv ${HOME}/blackbean.yaml ${HOME}/.blackbean.yaml

install: ## build blackbean.
	@go build -o bin/blackbean

clean: ## Remove test report.
	rm -f *.out
