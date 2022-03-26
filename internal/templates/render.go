package templates

import (
	"echo-starter/internal/models"
	"echo-starter/internal/wellknown"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"

	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, claimsPrincipal contracts_core_claimsprincipal.IClaimsPrincipal, code int, name string, data map[string]interface{}) error {
	data["isAuthenticated"] = func() bool { return claimsPrincipal.HasClaimType(wellknown.ClaimTypeAuthenticated) }
	data["paths"] = models.NewPaths()
	return c.Render(code, name, data)

}
