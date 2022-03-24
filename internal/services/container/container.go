package container

import (
	contracts_container "echo-starter/internal/contracts/container"
	"echo-starter/internal/shared"

	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/rs/zerolog/log"
)

// AddContainerAccessorFunc adds a function that returns a container.
func AddContainerAccessorFunc(builder *di.Builder) {
	log.Info().Str("DI", "ContainerAccessor").Send()
	contracts_container.AddContainerAccessorFunc(builder, shared.GetRootContainer)
}
