package client

import (
	"context"
	"time"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	conn  *grpc.ClientConn
	Auth  userpb.AuthServiceClient
	User  userpb.UserServiceClient
	Admin userpb.UserServiceClient // same service, different role checks
	Addr  userpb.AddressServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &UserClient{
		conn:  conn,
		Auth:  userpb.NewAuthServiceClient(conn),
		User:  userpb.NewUserServiceClient(conn),
		Admin: userpb.NewUserServiceClient(conn),
		Addr:  userpb.NewAddressServiceClient(conn),
	}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}
