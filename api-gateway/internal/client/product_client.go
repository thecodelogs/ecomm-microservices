package client

import (
	"context"
	"time"

	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn    *grpc.ClientConn
	Product productpb.ProductServiceClient
}

func NewProductClient(addr string) (*ProductClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &ProductClient{
		conn:    conn,
		Product: productpb.NewProductServiceClient(conn),
	}, nil
}

func (c *ProductClient) Close() error {
	return c.conn.Close()
}
