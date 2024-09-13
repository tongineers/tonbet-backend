package initializers

import (
	"fmt"

	tonlib "github.com/mercuryoio/tonlib-go/v2"
	"github.com/tongineers/tonlib-go-api/config"
	"github.com/tongineers/tonlib-go-api/internal/services/tonapi"
	"go.uber.org/zap"
)

func InitializeTonClient(conf *config.Config) (*tonlib.Client, error) {
	options, err := tonlib.ParseConfigFile(conf.TONLibConfigPath)
	if err != nil {
		return nil, fmt.Errorf("parse config file error: %w", err)
	}

	req := tonlib.TonInitRequest{
		Type:    "init",
		Options: *options,
	}

	client, err := tonlib.NewClient(&req, tonlib.Config{}, 10, true, 9)
	if err != nil {
		return nil, fmt.Errorf("creates client error: %w", err)
	}

	return client, nil
}

func InitializeTonClientOpts(
	сlient *tonlib.Client,
	conf *config.Config,
	logger *zap.Logger,
) []tonapi.Opt {
	return []tonapi.Opt{
		tonapi.WithClient(сlient),
		tonapi.WithConfig(conf),
		tonapi.WithLogger(logger),
	}
}
