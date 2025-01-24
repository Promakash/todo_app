PROTO_SRC_DIR := protos/proto
PROTO_GEN_DIR := protos/gen/go

.PHONY: proto_generate clean_proto

generate_docs:
	swag fmt
	swag init -g api-gateway/cmd/http/main.go -o api-gateway/docs

lint_code:
	golangci-lint run

proto_generate:
	protoc -I=$(PROTO_SRC_DIR) \
		$(PROTO_SRC_DIR)/*.proto \
		--go_out=$(PROTO_GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_GEN_DIR) --go-grpc_opt=paths=source_relative

clean_proto:
	rm -rf $(PROTO_GEN_DIR)/*
