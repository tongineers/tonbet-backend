package app

import (
	"context"

	"github.com/tongineers/tonbet-backend/internal/app/dependencies"
	"github.com/tongineers/tonbet-backend/internal/app/factories"
	"github.com/tongineers/tonbet-backend/pkg/workerpool"
)

// Application is a main struct for the application that contains general information
type Application struct {
	server    *factories.Server
	container *dependencies.Container
}

// InitializeApplication initializes new application
func InitializeApplication() (*Application, error) {

	return BuildApplication()
}

// Start starts application services
func (a *Application) Start(ctx context.Context, cli bool) error {
	if cli {
		return nil
	}

	a.container.Repository.MustMigrate()

	tasks := []workerpool.Task{
		a.container.Listener,
		a.container.Resolver,
		a.container.Fetcher,
	}

	pool := workerpool.NewWorkerPool(3)
	pool.Start()

	for _, task := range tasks {
		pool.Submit(task)
	}

	return a.server.Start(ctx)
}

// Stop stops application services
func (a *Application) Stop() (err error) {
	return a.server.Stop()
}
