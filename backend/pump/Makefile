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
		--proto_path=. \
		--go_out=./api/grpc/gen \
		--go_opt=module=github.com/jst-Frenzy/iot/backend/pump/api/grpc/gen \
		--go-grpc_out=./api/grpc/gen \
		--go-grpc_opt=module=github.com/jst-Frenzy/iot/backend/pump/api/grpc/gen \
		./api/grpc/proto/pump.proto
	@echo "[OK] Protobuf code generated successfully!"