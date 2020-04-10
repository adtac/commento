SHELL = bash

BUILD_DIR              = build
DEVEL_BUILD_DIR        = $(BUILD_DIR)/devel
PROD_BUILD_DIR         = $(BUILD_DIR)/prod

GO_SRC_DIR             = .
GO_SRC_FILES           = $(wildcard $(GO_SRC_DIR)/*.go)
GO_DEVEL_BUILD_DIR     = $(DEVEL_BUILD_DIR)
GO_DEVEL_BUILD_BINARY  = $(GO_DEVEL_BUILD_DIR)/commento
GO_PROD_BUILD_DIR      = $(PROD_BUILD_DIR)
GO_PROD_BUILD_BINARY   = $(GO_PROD_BUILD_DIR)/commento

devel: devel-go

prod: prod-go

test: test-go

clean:
	rm -rf $(BUILD_DIR)

# There's really no difference between the prod and devel binaries in Go, but
# for consistency sake, we'll use separate targets (maybe this will be useful
# later down the line).

devel-go:
	GO111MODULE=on go mod vendor
	GO111MODULE=on go build -mod=vendor -v -o $(GO_DEVEL_BUILD_BINARY) -ldflags "-X main.version=v1.8.0"

prod-go:
	GO111MODULE=on go mod vendor
	GO111MODULE=on go build -mod=vendor -v -o $(GO_PROD_BUILD_BINARY) -ldflags "-X main.version=v1.8.0"

test-go:
	GO111MODULE=on go mod vendor
	go test -v .

$(shell mkdir -p $(GO_DEVEL_BUILD_DIR) $(GO_PROD_BUILD_DIR))
