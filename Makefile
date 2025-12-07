PROTO_SRC = protos
PROTO_DEST = shared/protos

proto:
	protoc --proto_path=$(PROTO_SRC) \
		--go_out=$(PROTO_DEST) \
		--go-grpc_out=$(PROTO_DEST) \
		$(shell find $(PROTO_SRC) -name "*.proto")

.PHONY: proto
