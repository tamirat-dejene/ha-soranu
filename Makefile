PROTO_SRC = protos
PROTO_DEST = .
GO_MODULE = github.com/tamirat-dejene/ha-soranu

proto:
	protoc --proto_path=$(PROTO_SRC) \
		--go_out=$(PROTO_DEST) --go_opt=module=$(GO_MODULE) \
		--go-grpc_out=$(PROTO_DEST) --go-grpc_opt=module=$(GO_MODULE) \
		$(shell find $(PROTO_SRC) -name "*.proto")

.PHONY: proto
