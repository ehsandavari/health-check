serviceName=$(notdir $(CURDIR))
OsName=$(shell uname -s)

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"; printf "\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

init: install-deps swagger submodule proto ### init code base
.PHONY: init

run: linter-golangci swagger ### run
	go mod tidy && go mod download && \
	go run .
.PHONY: run

swagger: ### init swagger
	swag fmt
	swag init --parseDependency --parseInternal --parseDepth 1 --generalInfo ./presentation/api/v1/v1.go --output ./presentation/api/docs --instanceName v1 --generatedTime true
.PHONY: swagger

linter-golangci: ### check by golangci linter
	golangci-lint run ./...
.PHONY: linter-golangci

test: ### run unit-test and integration-test
	go test -v -cover -race ./... && \
    go clean -testcache && \
    go test -v ./test/...
.PHONY: test

install-deps: ### install all dependencies
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install go.uber.org/mock/mockgen@latest
	go install golang.org/x/tools/cmd/stringer@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
ifeq ($(OsName), Linux)
	sudo apt install -y protobuf-compiler
else ifeq ($(OsName), Darwin)
	brew install protobuf
else
	$(error "Unsupported operating system: $(OsName)")
endif
	protoc --version
.PHONY: install-deps

submodule: ### init submodule from git
	git submodule init
	git submodule update
	cd presentation/grpc/proto/ && git checkout develop
.PHONY: submodule

submodule-pull: ### update submodule in your project form git
	cd presentation/grpc/proto/$(serviceName) && \
	git pull
.PHONY: submodule-pull

submodule-push: ### update submodule in git | make proto-update m="your message"
	cd presentation/grpc/proto/$(serviceName) && \
	git add . && \
	git commit -m "$(m)" && \
	git push
.PHONY: submodule-push

submodule-status: ### submodule status in git
	cd presentation/grpc/proto/$(serviceName) && \
	git status
.PHONY: submodule-status

proto: ### init proto
	protoc --go_out=presentation/grpc/proto/$(serviceName)/ --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=presentation/grpc/proto/$(serviceName)/ presentation/grpc/proto/$(serviceName)/*.proto
.PHONY: proto
