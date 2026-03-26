package grpcclient

import (
	pb "calculator-service/api"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.StorageServiceClient
}

func New(addr string) *Client {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewCalculatorServiceClient(conn),
	}
}

func (c *Client) Save(expr string, result float64) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.client.SaveCalculation(ctx, &pb.SaveRequest{
		Expression: expr,
		Result:     result,
	})

	if err != nil {
		log.Println("gRPC error:", err)
	}
}

func (c *Client) Close() {
	c.conn.Close()
}
