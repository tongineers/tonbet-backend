package appgo

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type (
	App struct {
		httpServer *http.Server
		grpcServer *grpc.Server
		conf       *Config
		logger     *zap.Logger
	}

	GRPCService interface {
		ServiceDef() *GRPCOptions
	}

	GRPCOptions struct {
		Handler     func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
		ServiceDesc *grpc.ServiceDesc
		ServiceImpl any
	}
)

func New(conf *Config, logger *zap.Logger) *App {
	app := &App{
		conf:   conf,
		logger: logger,
	}

	if conf.EnableGrpc {
		app.grpcServer = grpc.NewServer()
	}

	app.httpServer = &http.Server{
		Addr: fmt.Sprintf("localhost:%d", conf.HttpPort),
	}

	return app
}

func (app *App) Start(ctx context.Context) error {
	if app.grpcServer != nil {
		addr := fmt.Sprintf(":%d", app.conf.GrpcPort)
		app.logger.Info("Started GRPC Server", zap.String("GRPCServerAddress", addr))

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			app.logger.Error("Failed to listen", zap.Error(err))
			return err
		}

		defer func() {
			if err := listener.Close(); err != nil {
				app.logger.Error("Failed to close TCP", zap.Error(err))
			}
		}()

		go func() {
			if err = app.grpcServer.Serve(listener); err != nil {
				app.logger.Error("GRPC Server stopped", zap.Error(err))
			}
		}()
	}

	app.logger.Info("Started HTTP Server", zap.String("HTTPServerAddress", app.httpServer.Addr))
	if err := app.httpServer.ListenAndServe(); err != nil {
		app.logger.Error("HTTP Server stopped", zap.Error(err))
	}

	return app.Stop()
}

func (app *App) RegisterServices(ctx context.Context, conn *grpc.ClientConn, services ...GRPCService) error {
	mux := runtime.NewServeMux()

	for _, service := range services {
		err := service.ServiceDef().Handler(ctx, mux, conn)
		if err != nil {
			return err
		}

		app.grpcServer.RegisterService(service.ServiceDef().ServiceDesc, service.ServiceDef().ServiceImpl)
	}

	app.httpServer.Handler = mux
	return nil
}

func (app *App) Stop() error {
	if app.grpcServer != nil {
		app.grpcServer.GracefulStop()
	}
	return app.httpServer.Shutdown(context.TODO())
}
