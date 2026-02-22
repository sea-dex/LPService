include version.mk

LDFLAGS += -X 'starbase.ag/liquidity/version.Version=$(VERSION)'
LDFLAGS += -X 'starbase.ag/liquidity/version.GitRev=$(GITREV)'
LDFLAGS += -X 'starbase.ag/liquidity/version.GitBranch=$(GITBRANCH)'
LDFLAGS += -X 'starbase.ag/liquidity/version.BuildDate=$(DATE)'

# Check dependencies
# Check for Go
.PHONY: check-go
check-go:
	@which go > /dev/null || (echo "Error: Go is not installed" && exit 1)

# Check for Docker
.PHONY: check-docker
check-docker:
	@which docker > /dev/null || (echo "Error: docker is not installed" && exit 1)

# Check for Curl
.PHONY: check-curl
check-curl:
	@which curl > /dev/null || (echo "Error: curl is not installed" && exit 1)

# Targets that require the checks
build: check-go
lint: check-go
build-docker: check-docker
build-docker-nc: check-docker

.PHONY: build
build: ## Builds the binary locally into ./dist
# go build -ldflags "all=$(LDFLAGS)" -o liquidity ./cmd/liquid
# go build -ldflags "all=$(LDFLAGS)" -o eventfeed ./cmd/eventd
# go get -tags musl -u github.com/confluentinc/confluent-kafka-go/v2@v2.5.3
	go build -tags musl -ldflags "all=$(LDFLAGS)" -o lpservice ./cmd/LPTracker

.PHONY: build-docker
build-docker: ## Builds a docker image with the node binary
	docker build -t zkevm-node -f ./Dockerfile .

.PHONY: build-docker-nc
build-docker-nc: ## Builds a docker image with the node binary - but without build cache
	docker build --no-cache=true -t zkevm-node -f ./Dockerfile .

.PHONY: dep
dep: ## Install dep tool
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install mvdan.cc/gofumpt@latest
	go install github.com/daixiang0/gci@latest
	go install github.com/4meepo/tagalign/cmd/tagalign@latest
	go install github.com/bombsimon/wsl/v4/cmd/wsl@master
	go install github.com/tetafro/godot/cmd/godot@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.61.0

.PHONY: lint
lint: ## Runs the linter
	golangci-lint run ./...

.PHONY: abi
abi: ## generate abi go file
	abigen --abi contracts/abi/erc20.json --pkg contracts --type ERC20 --out contracts/erc20.go
	abigen --abi contracts/abi/erc721.json --pkg contracts --type ERC721 --out contracts/erc721.go
	abigen --abi contracts/abi/erc1155.json --pkg contracts --type ERC1155 --out contracts/erc1155.go
	abigen --abi contracts/abi/multicall3.json --pkg contracts --type Multicall3 --out contracts/multicall3.go
	abigen --abi contracts/abi/StarBaseLimitOrder.json --pkg limitorder --type LimitOrder --out contracts/limitorder/limitorder.go
	abigen --abi contracts/abi/StarBaseLimitOrderBot.json --pkg limitorder --type LimitOrderBot --out contracts/limitorder/limitorderbot.go
	abigen --abi contracts/abi/StarBaseDCA.json --pkg dca --type DCA --out contracts/dca/dca.go
	abigen --abi contracts/abi/StarBaseDCABot.json --pkg dca --type DCAbot --out contracts/dca/dcabot.go
	abigen --abi contracts/abi/UniswapV3Pool.json --pkg swaprouter --type UniswapV3Pool --out contracts/swaprouter/UniswapV3Pool.go
	abigen --abi contracts/abi/UniswapV2Pair.json --pkg swaprouter --type UniswapV2Pair --out contracts/swaprouter/UniswapV2Pair.go
	abigen --abi contracts/abi/UniswapV3Factory.json --pkg swaprouter --type UniswapV3Factory --out contracts/swaprouter/UniswapV3Factory.go
	abigen --abi contracts/abi/AggregatedSwapRouter.json --pkg swaprouter --type AggregatedSwapRouter --out contracts/swaprouter/AggregatedSwapRouter.go
	abigen --abi contracts/abi/SwapRouterMock.json --pkg swaproutermock --type swaproutermock --out contracts/swaproutermock/SwapRouterMock.go
	abigen --abi contracts/abi/UniswapV3PoolQuery.json --pkg poolquery --type UniswapV3PoolQuery --out contracts/poolquery/UniswapV3PoolQuery.go
	abigen --abi contracts/abi/PancakeV3Pool.json --pkg swaprouter --type PancakeV3Pool --out contracts/swaprouter/PancakeV3Pool.go
	abigen --abi contracts/abi/AeroV2Pool.json --pkg swaprouter --type AeroV2Pool --out contracts/swaprouter/AeroV2Pool.go
	abigen --abi contracts/abi/ArbitrageContract.json --pkg arbitrage --type ArbitrageContract --out contracts/arbitrage/ArbitrageContract.go

.PHONY: test
test: ## go test
	go test  -short -v ./...

.PHONY: fmt
fmt: ## go fmt
	go fmt ./...
	gofumpt -l -w .
	gci write  -s standard -s default -s "prefix(periphery)" .
	swag fmt

.PHONY: fix
fix: fmt ## auto fix code
	wsl --fix ./...
	tagalign -fix -sort ./...
	godot -w ./

.PHONY: clean
clean: ## clear cache
	rm -rf liquidity eventfeed lpservice lpservice.log

update-package: ## update all go package
	go get -u ./...
	go mod tidy -v

.PHONY: docs
docs:
	swag init -g api.go -d ./api  --parseDependency

.PHONY: ci
ci: lint fmt fix build test
	echo $? && echo "success!"

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := build

.PHONY: help
help: ## Prints this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
