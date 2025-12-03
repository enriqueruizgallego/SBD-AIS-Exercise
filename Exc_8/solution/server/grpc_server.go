package server

import (
	"context"
	"exc8/pb"
	"net"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPCService struct {
	pb.UnimplementedOrderServiceServer
}

func StartGrpcServer() error {
	// Create a new gRPC server.
	srv := grpc.NewServer()
	// Create grpc service
	grpcService := &GRPCService{}
	// Register our service implementation with the gRPC server.
	pb.RegisterOrderServiceServer(srv, grpcService)
	// Serve gRPC server on port 4000.
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		return err
	}
	err = srv.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

var drinks = []*pb.Drink{
	{Id: 1, Name: "Spritzer", Price: 2, Description: "Wine with soda"},
	{Id: 2, Name: "Beer", Price: 3, Description: "Hagenberger Gold"},
	{Id: 3, Name: "Coffee", Description: "Mifare isn't that secure"},
}

var orders = map[int32]int32{}

// todo implement functions

func (s *GRPCService) GetDrinks(ctx context.Context, _ *emptypb.Empty) (*pb.DrinkList, error) {
	return &pb.DrinkList{
		Drinks: drinks,
	}, nil
}

func (s *GRPCService) OrderDrink(ctx context.Context, req *pb.OrderItem) (*wrapperspb.BoolValue, error) {
	//check if the drink exist
	var exists bool
	for _, d := range drinks {
		if d.Id == req.DrinkId {
			exists = true
			break
		}
	}
	// if drink not exist,we return false
	if !exists {
		return wrapperspb.Bool(false), nil
	}
	//if exist, add the quantity to the orders
	orders[req.DrinkId] += req.Quantity
	//return true to tell that all is okey
	return wrapperspb.Bool(true), nil
}

func (s *GRPCService) GetOrders(ctx context.Context, _ *emptypb.Empty) (*pb.OrderSummary, error) {
	summary := &pb.OrderSummary{}
	//see all the orders
	for drinkID, totalQty := range orders {
		// find the drink associated with that id
		var drink *pb.Drink
		for _, d := range drinks {
			if d.Id == drinkID {
				drink = d
				break
			}
		}
		// create a OrderSummaryItem with the drink an the total quantity
		item := &pb.OrderSummaryItem{
			Drink:         drink,
			TotalQuantity: totalQty,
		}
		//add to summary itemlist
		summary.Items = append(summary.Items, item)
	}
	//return the total summary
	return summary, nil
}
