package handler

import (
	"context"
	"log"
	"strings"

	"github.com/Fulim13/grpc-demo/gen/golang/v1/order"
	"github.com/Fulim13/grpc-demo/gen/golang/v1/payment"
	"github.com/Fulim13/grpc-demo/services/order/internal/client"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	order.UnimplementedOrderServer
	// service service.OrderService
}

// func NewGRPCHandler(server *grpc.Server, service service.OrderService) *grpcHandler {
func NewGRPCHandler(server *grpc.Server) *grpcHandler {
	handler := &grpcHandler{
		// service: service,
	}
	order.RegisterOrderServer(server, handler)
	reflection.Register(server) // exposes the server's service/method list so clients like grpcurl can call it without a .proto file
	return handler
}

func (h *grpcHandler) Create(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	userId := req.GetUserId()
	orderItems := req.GetOrderItems()

	log.Printf("Receive userID: %v\n", userId)
	log.Printf("Receive orderItems: %+v\n", orderItems)

	// if userId < 1 {
	// 	return nil, status.Errorf(codes.Internal, "userId must be positive number")
	// }

	// h.service.CreateOrder(userId, orderItem)

	// Call to payment service
	paymentService, err := client.NewPaymentServiceClient()
	if err != nil {
		return nil, err
	}

	// Timeout Pattern
	// ctx, _ := context.WithTimeout(context.TODO(), time.Second*3)

	res, paymentErr := paymentService.Client.Create(context.Background(), &payment.CreatePaymentRequest{
		UserId:     userId,
		OrderId:    1,
		TotalPrice: 10,
	})

	log.Printf("Receive response: %v\n", res)

	if paymentErr != nil {
		st := status.Convert(paymentErr)
		var allErrors []string
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case *errdetails.BadRequest:
				for _, violation := range t.GetFieldViolations() {
					allErrors = append(allErrors, violation.Description)
				}
			}
		}
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: strings.Join(allErrors, "\n"),
		}
		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")
		statusWithDetails, _ := orderStatus.WithDetails(badReq)
		return nil, statusWithDetails.Err()
	}
	return &order.CreateOrderResponse{OrderId: 1}, nil
}
