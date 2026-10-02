package handler

import (
	"context"

	"github.com/Fulim13/grpc-demo/gen/golang/v1/payment"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	payment.UnimplementedPaymentServer
	// service service.PaymentService
}

// func NewGRPCHandler(server *grpc.Server, service service.PaymentService) *grpcHandler {
func NewGRPCHandler(server *grpc.Server) *grpcHandler {
	handler := &grpcHandler{
		// service: service,
	}
	payment.RegisterPaymentServer(server, handler)
	reflection.Register(server) // exposes the server's service/method list so clients like grpcurl can call it without a .proto file
	return handler
}

func (h *grpcHandler) Create(ctx context.Context, req *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	//Timeout Pattern
	// time.Sleep(4 * time.Second)

	// Retry and Circuilt Breaker Pattern
	// fmt.Println("In Payment service")
	// return nil, status.New(codes.Unavailable, fmt.Sprint("payment service unavailable.")).Err()

	var validationErrors []*errdetails.BadRequest_FieldViolation

	if req.UserId < 1 {
		validationErrors = append(validationErrors, &errdetails.BadRequest_FieldViolation{
			Field:       "user_id",
			Description: "user_id cannot be less than 1",
		})
	}

	if req.OrderId < 1 {
		validationErrors = append(validationErrors, &errdetails.BadRequest_FieldViolation{
			Field:       "order_id",
			Description: "order_id cannot be less than 1",
		})
	}

	if req.TotalPrice < 0 {
		validationErrors = append(validationErrors, &errdetails.BadRequest_FieldViolation{
			Field:       "total_price",
			Description: "total_price cannot be less than 0",
		})
	}

	if len(validationErrors) > 0 {
		stat := status.New(400, "invalid payment request")
		badRequest := &errdetails.BadRequest{}
		badRequest.FieldViolations = validationErrors
		s, _ := stat.WithDetails(badRequest)
		return nil, s.Err()
	}

	// payment, err := h.service.CreatePayment(req.userId, req.orderId, req.TotalPrice)
	// if err != nil {
	// 	return nil, status.New(codes.Internal, fmt.Sprintf("failed to charge. %v ", err)).Err()
	// }
	return &payment.CreatePaymentResponse{PaymentId: 1}, nil
}
