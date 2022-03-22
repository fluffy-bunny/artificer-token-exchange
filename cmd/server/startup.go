package main

import (
	tex_config "echo-starter/internal/contracts/config"

	core_contracts "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/core"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Startup struct {
	config *tex_config.Config
}

func NewStartup() *Startup {
	return &Startup{
		config: &tex_config.Config{},
	}
}
func (s *Startup) GetPort() int {
	return s.config.Port
}
func (s *Startup) GetConfigOptions() *core_contracts.ConfigOptions {
	return &core_contracts.ConfigOptions{
		RootConfig:  []byte(tex_config.ConfigDefaultJSON),
		Destination: s.config,
	}
}
func (s *Startup) ConfigureServices(builder *di.Builder) error {
	di.AddSingletonTypeByObj(builder, s.config)

	return nil
}
func (s *Startup) Configure(e *echo.Echo, root di.Container) error {
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			id := uuid.New()
			return id.String()
		},
	}))
	return nil
}
