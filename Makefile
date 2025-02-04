PROTO_ROOT := .
PROTO_OUT := libs/pb

.PHONY: generate
generate:
	# Генерация для customers
	protoc -I=$(PROTO_ROOT)/libs/protobuf \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_ROOT)/libs/protobuf/customers/customer.proto
	
	# Генерация для shared
	protoc -I=$(PROTO_ROOT)/libs/protobuf \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		$(PROTO_ROOT)/libs/protobuf/shared/common.proto
