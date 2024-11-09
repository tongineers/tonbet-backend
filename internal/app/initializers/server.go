package initializers

import (
	"context"
	"fmt"

	"github.com/tongineers/dice-ton-api/internal/app/dependencies"
	appgo "github.com/tongineers/dice-ton-api/pkg/app-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitializeServer(container *dependencies.Container) (*appgo.App, error) {
	app := appgo.New(&appgo.Config{
		GrpcPort:   container.Config.AppPort,
		HttpPort:   container.Config.AppHttpPort,
		EnableGrpc: true,
	}, container.Logger)

	conn, err := dialTCP(fmt.Sprintf(":%d", container.Config.AppPort))
	if err != nil {
		return nil, err
	}

	return app, app.RegisterServices(context.Background(), conn, container.Service)
}

func dialTCP(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
