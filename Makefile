LOCAL_BIN := $(CURDIR)/bin
PROTO_DIR := $(CURDIR)/proto
EXAMPLE_DIR := $(CURDIR)/example
VENDOR_PROTO_DIR := $(CURDIR)/vendor.protogen
MODULE := github.com/casnerano/protoc-gen-go-rbac

.PHONY: download-bin-deps
download-bin-deps:
	ls $(LOCAL_BIN)/buf &> /dev/null || GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@latest
	ls $(LOCAL_BIN)/protoc-gen-go &> /dev/null || GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	ls $(LOCAL_BIN)/protoc-gen-go-grpc &> /dev/null || GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

.PHONY: build-protoc-gen-go-rbac
build-protoc-gen-go-rbac:
	go build -o ${LOCAL_BIN}/protoc-gen-go-rbac ./cmd/protoc-gen-go-rbac

.PHONY: generate
generate: download-bin-deps
	$(LOCAL_BIN)/buf generate
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(LOCAL_BIN)
	rm -rf $(CURDIR)/proto/*.pb.go
	rm -rf $(CURDIR)/example/pb/*.pb.go

.PHONY: vendor.protogen/google/api
vendor.protogen/google/api:
	@if [ ! -d "$(VENDOR_PROTO_DIR)/google/api" ]; then \
		echo "Vendoring google/api from googleapis..."; \
		git clone -b master --single-branch -n --depth=1 --filter=tree:0 https://github.com/googleapis/googleapis "$(VENDOR_PROTO_DIR)/googleapis" && \
		cd "$(VENDOR_PROTO_DIR)/googleapis" && git sparse-checkout set --no-cone google/api && git checkout; \
		mkdir -p "$(VENDOR_PROTO_DIR)/google"; \
		mv "$(VENDOR_PROTO_DIR)/googleapis/google/api" "$(VENDOR_PROTO_DIR)/google/"; \
		rm -rf "$(VENDOR_PROTO_DIR)/googleapis"; \
	fi

.PHONY: vendor.protogen
vendor.protogen: vendor.protogen/google/api