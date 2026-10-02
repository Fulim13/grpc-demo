## Generate protobuf code

````sh
# 1. Create proto/order.proto, proto/payment.proto

# 2. Install Required Package for grpc
```sh
brew install protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
brew install grpcurl
````

# 3. Generate the protobuf code

```sh
protoc -I proto \
  --go_out=. --go_opt=module=github.com/Fulim13/grpc-demo \
  --go-grpc_out=. --go-grpc_opt=module=github.com/Fulim13/grpc-demo \
  proto/order.proto proto/payment.proto
```

# 4. Initialize gen/golang/v1/order and gen/golang/v1/payment module

```sh
cd gen/golang/v1/order
go mod init github.com/Fulim13/grpc-demo/gen/golang/v1/order
go mod tidy

cd ../payment
go mod init github.com/Fulim13/grpc-demo/gen/golang/v1/payment
go mod tidy
```

# 5. Initialize services/order + services/payment modules, and link them to gen/ with go.work

```
  services/order  ──require──>  gen/golang/v1/order    (generated pb.go)
                  ──require──>  gen/golang/v1/payment
  services/payment ─require──>  gen/golang/v1/payment

  replace + go.work  =>  resolve to the LOCAL folders, not the network
                         => edit gen/*.pb.go and services/ pick it up instantly
```

```sh
# services/order
cd services/order
go mod init github.com/Fulim13/grpc-demo/services/order
go mod edit -replace github.com/Fulim13/grpc-demo/gen/golang/v1/order=../../gen/golang/v1/order
go mod edit -replace github.com/Fulim13/grpc-demo/gen/golang/v1/payment=../../gen/golang/v1/payment
go mod tidy

# services/payment
cd ../payment
go mod init github.com/Fulim13/grpc-demo/services/payment
go mod edit -replace github.com/Fulim13/grpc-demo/gen/golang/v1/payment=../../gen/golang/v1/payment
go mod tidy

# workspace at repo root
cd ../..
go work init ./gen/golang/v1/order ./gen/golang/v1/payment ./services/order ./services/payment
go work sync
```

## How to Run

```sh
cd services/payment
APPLICATION_PORT=1324 go run cmd/server/main.go

cd services/order
APPLICATION_PORT=1323 PAYMENT_SERVICE_URL=localhost:1324 go run cmd/server/main.go

grpcurl \
	-d '{"user_id": 123, "order_items": [{"product_code": "prod", "quantity": 4, "unit_price": 12}]}' \
	-plaintext \
	localhost:1323 \
	order.v1.Order/Create

# send values the payment service rejects
grpcurl \
	-d '{"user_id": -1, "order_items": [{"product_code": "prod", "quantity": 4, "unit_price": 12}]}' \
	-plaintext \
	localhost:1323 \
	order.v1.Order/Create
```
