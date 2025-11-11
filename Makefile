PROTO_SRC = proto
PROTO_DEST = shared/proto

proto:
	protoc --proto_path=$(PROTO_SRC) \
		--go_out=$(PROTO_DEST) \
		--go-grpc_out=$(PROTO_DEST) \
		$(shell find $(PROTO_SRC) -name "*.proto")

.PHONY: proto
