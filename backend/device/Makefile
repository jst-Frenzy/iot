.PHONY:  \
    lint \
    proto-gen \


lint:
	@echo "[INFO] Running golangci-lint..."
	@golangci-lint run
	@echo "[OK] Lint completed successfully!"

proto-gen:
	@echo "[INFO] Generating protobuf code..."
	@protoc --version
	@protoc \
		--experimental_allow_proto3_optional \
		--proto_path= \
		--go_out=./api/grpc/gen \
		--go_opt=module=gl.iteco.com/technology/go_masterdata/preorderservice/api/grpc/gen \
		--go-grpc_out=./api/grpc/gen \
		--go-grpc_opt=module=gl.iteco.com/technology/go_masterdata/preorderservice/api/grpc/gen \
		./api/grpc/proto/preorders/preorders.proto \
		./api/grpc/proto/offers/offers.proto \
		./api/grpc/proto/trades/trades.proto
	@echo "[OK] Protobuf code generated successfully!"