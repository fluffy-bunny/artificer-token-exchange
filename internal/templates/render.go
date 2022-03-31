package templates

import (
	"echo-starter/internal/models"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	core_echo_templates "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/templates"
	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, claimsPrincipal contracts_core_claimsprincipal.IClaimsPrincipal, code int, name string, data map[string]interface{}) error {
	data["isAuthenticated"] = func() bool {
		return claimsPrincipal.HasClaimType(core_wellknown.ClaimTypeAuthenticated)
	}
	data["paths"] = models.NewPaths()
	data["claims"] = claimsPrincipal.GetClaims()
	return core_echo_templates.Render(c, code, name, data)

}
