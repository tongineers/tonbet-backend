package factories

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type (
	Server struct {
		httpServer *http.Server
		grpcServer *grpc.Server
		config     *ServerConfig
		logger     *zap.Logger
	}

	ServerConfig struct {
		HttpPort   int
		GrpcPort   int
		Router     *gin.Engine
		EnableGrpc bool
	}

	Service interface {
		ServiceDef() *Options
	}

	Options struct {
		Handler     func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
		ServiceDesc *grpc.ServiceDesc
		ServiceImpl any
	}
)

func ServerFactory(config *ServerConfig, logger *zap.Logger) *Server {
	srv := &Server{
		config: config,
		logger: logger,
	}

	if config.EnableGrpc {
		srv.grpcServer = grpc.NewServer()
	}

	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", config.HttpPort),
		Handler: config.Router,
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	if s.grpcServer != nil {
		addr := fmt.Sprintf(":%d", s.config.GrpcPort)
		s.logger.Info("Started GRPC Server", zap.String("GRPCServerAddress", addr))

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Error("Failed to listen", zap.Error(err))
			return err
		}

		defer func() {
			if err := listener.Close(); err != nil {
				s.logger.Error("Failed to close TCP", zap.Error(err))
			}
		}()

		go func() {
			if err = s.grpcServer.Serve(listener); err != nil {
				s.logger.Error("GRPC Server stopped", zap.Error(err))
			}
		}()
	}

	s.logger.Info("Started HTTP Server", zap.String("HTTPServerAddress", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.logger.Error("HTTP Server stopped", zap.Error(err))
	}

	return s.Stop()
}

func (s *Server) RegisterServices(ctx context.Context, conn *grpc.ClientConn, services ...Service) error {
	mux := runtime.NewServeMux()

	for _, service := range services {
		err := service.ServiceDef().Handler(ctx, mux, conn)
		if err != nil {
			return err
		}

		s.grpcServer.RegisterService(service.ServiceDef().ServiceDesc, service.ServiceDef().ServiceImpl)
	}

	s.httpServer.Handler = mux
	return nil
}

func (s *Server) Stop() error {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	return s.httpServer.Shutdown(context.TODO())
}
