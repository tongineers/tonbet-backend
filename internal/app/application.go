package app

import (
	"context"

	"github.com/tongineers/dice-ton-api/internal/app/dependencies"
	appgo "github.com/tongineers/dice-ton-api/pkg/app-go"
)

// Application is a main struct for the application that contains general information
type Application struct {
	app       *appgo.App
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

	return a.app.Start(ctx)
}

// Stop stops application services
func (a *Application) Stop() error {
	return a.app.Stop()
}
