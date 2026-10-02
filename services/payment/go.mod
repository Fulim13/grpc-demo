module github.com/Fulim13/grpc-demo/services/payment

go 1.27.0

replace github.com/Fulim13/grpc-demo/gen/golang/v1/order => ../../gen/golang/v1/order

replace github.com/Fulim13/grpc-demo/gen/golang/v1/payment => ../../gen/golang/v1/payment

require (
	github.com/Fulim13/grpc-demo/gen/golang/v1/payment v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260928230214-8a89bd6388cc
	google.golang.org/grpc v1.84.0
)

require (
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)
