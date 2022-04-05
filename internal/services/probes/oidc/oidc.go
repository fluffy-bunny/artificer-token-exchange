package oidc

import (
	"reflect"

	contracts_probe "echo-starter/internal/contracts/probe"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	di "github.com/fluffy-bunny/sarulabsdi"
)

type (
	service struct {
		Logger contracts_logger.ILogger `inject:""`
	}
)

func assertImplementation() {
	var _ contracts_probe.IProbe = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddSingletonIProbe registers the *service as a singleton.
func AddSingletonIProbe(builder *di.Builder) {
	contracts_probe.AddSingletonIProbe(builder, reflectType)
}
func (s *service) GetName() string {
	return "oidc"
}
func (s *service) Probe() error {
	s.Logger.Debug().Str("probe", "oidc").Send()

	return nil
}
