package templates

import (
	"echo-starter/internal/models"
	"echo-starter/internal/wellknown"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"

	"github.com/labstack/echo/v4"
)

func GetTemplateRender(rootDir string) *TemplateRenderer {
	t, err := FindAndParseTemplates(rootDir, nil)
	if err != nil {
		panic(err)
	}

	return &TemplateRenderer{
		templates: t,
	}
}

func Render(c echo.Context, claimsPrincipal contracts_core_claimsprincipal.IClaimsPrincipal, code int, name string, data map[string]interface{}) error {
	data["isAuthenticated"] = func() bool { return claimsPrincipal.HasClaimType(wellknown.ClaimTypeAuthenticated) }
	data["paths"] = models.NewPaths()
	return c.Render(code, name, data)

}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func FindAndParseTemplates(rootDir string, funcMap template.FuncMap) (*template.Template, error) {
	cleanRoot := filepath.Clean(rootDir)
	pfx := len(cleanRoot) + 1
	root := template.New("")

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".tpl") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(funcMap)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}

		return nil
	})

	return root, err
}
