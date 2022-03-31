package ready

import (
	"context"
	"echo-starter/internal/wellknown"
	"net/http"
	"reflect"
	"time"

	contracts_probe "echo-starter/internal/contracts/probe"

	"github.com/catmullet/go-workers"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	"github.com/rs/zerolog/log"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger contracts_logger.ILogger `inject:""`
		Probes []contracts_probe.IProbe `inject:""`
		runner workers.Runner
	}
	probeWorker struct{}
)

func assertImplementation() {
	var _ contracts_handler.IHandler = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder *di.Builder) {
	contracts_handler.AddScopedIHandlerEx(builder,
		reflectType,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
		},
		wellknown.ReadyPath)
}

func (s *service) Ctor() {
	s.runner = workers.NewRunner(context.Background(), &probeWorker{}, int64(len(s.Probes))).Start()
	s.runner.SetTimeout(time.Second * 5)
}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {
	for _, probe := range s.Probes {
		s.Logger.Debug().Str("probe", probe.GetName()).Msg("issuing probe")
		s.runner.Send(probe)
	}
	s.Logger.Debug().Msg("Waiting for probes to complete")
	err := s.runner.Wait()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "ok")
}
func (s *probeWorker) Work(in interface{}, out chan<- interface{}) error {
	var err error
	if probe, ok := in.(contracts_probe.IProbe); ok {
		log.Debug().Str("probe", probe.GetName()).Msg("probe it")
		err = probe.Probe()
		if err != nil {
			log.Error().Err(err).Msg("probe failed")
		}
	}
	return err

}
