GOMODULE := $(shell go list -m)
BIN_DIR  := $(CURDIR)/bin
export PATH := $(BIN_DIR):$(PATH)


BUILD_DIR := $(CURDIR)/build
.PHONY: build
build: dep
	go build -o $(BUILD_DIR)/server/server $(CURDIR)/cmd

.PHONY: dep
dep:
	go mod download

.PHONY: test
test:
	go test -v -cover -coverpkg ./internal/... ./internal/...
 
.PHONY: lint
lint:
	golangci-lint run ./cmd/... ./internal/... --allow-parallel-runners;

.PHONY: doc
doc:
	swag init -g cmd/main.go -o generate/doc;

PROTO_DIR           := $(CURDIR)/protobufs
PROTO_SOURCE_DIR    := $(CURDIR)/generate/pb
PROTO_IMPORT_PREFIX := $(GOMODULE)/generate/pb
.PHONY: proto
proto:
	rm -rf $(PROTO_SOURCE_DIR)
	protoc -I $(PROTO_DIR) \
	--go_out . --go-grpc_out=require_unimplemented_servers=false:. \
	--go_opt module=$(GOMODULE) --go-grpc_opt module=$(GOMODULE) \
	--go_opt Muser/user.proto=$(PROTO_IMPORT_PREFIX)/user \
	--go-grpc_opt Muser/user.proto=$(PROTO_IMPORT_PREFIX)/user \
	$(PROTO_DIR)/user/*

MOCK_SOURCE_DIR := $(CURDIR)/generate/mock
.PHONY: mock
mock:
	rm -rf $(MOCK_SOURCE_DIR)/*_mock.go
	mockgen -package mock -source ./internal/repository/user.go -destination $(MOCK_SOURCE_DIR)/user_mock.go
	mockgen -package mock -source ./internal/repository/user_token.go -destination $(MOCK_SOURCE_DIR)/user_token_mock.go
	mockgen -package mock -source ./internal/repository/fanclub.go -destination $(MOCK_SOURCE_DIR)/fanclub_mock.go
	mockgen -package mock -source ./internal/gim/gim_client.go -destination $(MOCK_SOURCE_DIR)/gim_client_mock.go
	mockgen -package mock -source ./internal/uuid/uuid.go -destination $(MOCK_SOURCE_DIR)/uuid_mock.go

.PHONY: install-bin uninstall-bin
install-bin:   | $(BIN_DIR)/swag $(BIN_DIR)/mockgen $(BIN_DIR)/golangci-lint $(BIN_DIR)/protoc-gen-go $(BIN_DIR)/protoc-gen-go-grpc $(BIN_DIR)/protoc
uninstall-bin: ; rm -rf $(BIN_DIR)

.PHONY: fmt-schema
fmt-schema:
	-atlas schema inspect \
		--env script --format "{{ sql . \" \" }}" \
		> schema.tmp && cp schema.tmp script/init.sql
	rm schema.tmp

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(BIN_DIR)/%: | $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install $(PACKAGE)

$(BIN_DIR)/swag:               PACKAGE=github.com/swaggo/swag/cmd/swag@v1.8.12
$(BIN_DIR)/mockgen:            PACKAGE=github.com/golang/mock/mockgen@v1.6.0
$(BIN_DIR)/golangci-lint:      PACKAGE=github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
$(BIN_DIR)/protoc-gen-go:      PACKAGE=google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
$(BIN_DIR)/protoc-gen-go-grpc: PACKAGE=google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0


ifeq (${shell uname}, Darwin)
	PROTOC_OS = osx
else ifeq (${shell uname}, Linux)
	PROTOC_OS = linux
else
	${error protobuf compiler unsupported current os}
endif
PROTOC_VERSION      := 21.1
PROTOC_DOWNLOAD_URI := https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-$(PROTOC_OS)-x86_64.zip
$(BIN_DIR)/protoc:
	curl -sL $(PROTOC_DOWNLOAD_URI) -o /tmp/protoc.zip && \
	unzip -oXq /tmp/protoc.zip -d /tmp/protoc && \
	mv -f /tmp/protoc/bin/protoc $(BIN_DIR)/protoc && \
	mv -f /tmp/protoc/include/ $(BIN_DIR)/include/ && \
	rm -rf /tmp/protoc.zip /tmp/protoc/
