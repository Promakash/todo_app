PROTO_SRC_DIR := protos/proto
PROTO_GEN_DIR := protos/gen/go

proto_generate:
	protoc -I=$(PROTO_SRC_DIR) \
		$(PROTO_SRC_DIR)/*.proto \
		--go_out=$(PROTO_GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR) --go-grpc_opt=paths=source_relative

clean_proto:
	rm -rf $(PROTO_GEN_DIR)/*

.PHONY: proto_generate clean_proto
