package client

import (
	"context"
	"exc8/pb"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcClient struct {
	client pb.OrderServiceClient
}

func NewGrpcClient() (*GrpcClient, error) {
	conn, err := grpc.NewClient(":4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewOrderServiceClient(conn)
	return &GrpcClient{client: client}, nil
}

func (c *GrpcClient) Run() error {
	// todo
	ctx := context.Background()
	// 1. List drinks
	fmt.Println("Requesting drinks 🍹🍺☕")

	drinksResp, err := c.client.GetDrinks(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("GetDrinks failed: %w", err)
	}

	fmt.Println("Available drinks:")
	for _, d := range drinksResp.Drinks {
		fmt.Printf("\t> id:%d  name:%q  price:%d  description:%q\n",
			d.Id, d.Name, d.Price, d.Description)
	}

	// 2. Order a few drinks
	fmt.Println("Ordering drinks 👨‍🍳⏱️🍻🍻")
	for _, d := range drinksResp.Drinks {
		fmt.Printf("\t> Ordering: 2 x %s\n", d.Name)

		_, err := c.client.OrderDrink(ctx, &pb.OrderItem{
			DrinkId:  d.Id,
			Quantity: 2,
		})
		if err != nil {
			return fmt.Errorf("OrderDrink failed: %w", err)
		}
	}
	// 3. Order more drinks
	fmt.Println("Ordering another round of drinks 👨‍🍳⏱️🍻🍻")
	for _, d := range drinksResp.Drinks {
		fmt.Printf("\t> Ordering: 6 x %s\n", d.Name)

		_, err := c.client.OrderDrink(ctx, &pb.OrderItem{
			DrinkId:  d.Id,
			Quantity: 6,
		})
		if err != nil {
			return fmt.Errorf("OrderDrink failed (second round): %w", err)
		}
	}
	// 4. Get order total
	fmt.Println("Getting the bill 💹💹💹")

	summary, err := c.client.GetOrders(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("GetOrders failed: %w", err)
	}

	for _, item := range summary.Items {
		fmt.Printf("\t> Total: %d x %s\n", item.TotalQuantity, item.Drink.Name)
	}

	return nil
}
