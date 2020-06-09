
VERSION := $(shell echo $(shell git describe --always) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')
all: ci-lint install

###############################################################################
# Build / Install
###############################################################################

LD_FLAGS = -X gitlab.bianjie.ai/irishub/hack-goz/cmd.Version=$(VERSION) \
	-X gitlab.bianjie.ai/irishub/hack-goz/cmd.Commit=$(COMMIT) \
	-X github.com/cosmos/cosmos-sdk/types.reDnmString=[a-z][a-z0-9/:-]{2,63}

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building hackcli binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/hackcli.exe main.go
else
	@echo "building hackcli binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/hackcli main.go
endif

install: go.sum
	@echo "installing hackcli binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o $${GOBIN-$${GOPATH-$$HOME/go}/bin}/hackcli main.go