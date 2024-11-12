package initializers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/envy"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gw "github.com/tongineers/dice-ton-api/gen/go/tonapi/v1"
)

const (
	// HTTPServerAddrEnv is an environment variable name for HTTP server address
	HTTPServerAddrEnv = "HTTP_SERVER_ADDR"
	// DefaultHTTPServerAddr  is a default value for HTTP server address
	DefaultHTTPServerAddr = ":8000"
)

// HTTPServerAddr is a type alias for HTTP server address values
type HTTPServerAddr string

// HTTPServerConfig is a type to store HTTP server config
type HTTPServerConfig struct {
	HTTPServerAddr HTTPServerAddr
	Router         *gin.Engine
}

// InitializeHTTPServerConfig initializes new config for InitializeHTTPServer
func InitializeHTTPServerConfig(router *gin.Engine) *HTTPServerConfig {
	return &HTTPServerConfig{
		HTTPServerAddr: HTTPServerAddr(envy.Get(HTTPServerAddrEnv, DefaultHTTPServerAddr)),
		Router:         router,
	}
}

// InitializeHTTPServer create new http.Server instance
func InitializeHTTPGateway(cfg *HTTPServerConfig) (*http.Server, error) {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterTonApiServiceHandlerFromEndpoint(context.Background(), mux, "localhost:9090", opts)
	if err != nil {
		return nil, err
	}

	// create http server
	srv := &http.Server{
		Addr:    string(cfg.HTTPServerAddr),
		Handler: mux,
	}

	return srv, nil
}
