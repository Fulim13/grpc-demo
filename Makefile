MODULE := github.com/Fulim13/grpc-demo

.PHONY: proto build clean

# Generate Go code from proto/*.proto
proto:
	protoc -I proto \
	  --go_out=. --go_opt=module=$(MODULE) \
	  --go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
	  proto/*.proto

# Delete generated code
clean:
	find gen -name '*.pb.go' -delete
