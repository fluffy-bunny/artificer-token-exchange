package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"echo-starter/internal/startup"

	echo_contracts_startup "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/startup"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/runtime"
	"github.com/golang/mock/gomock"
)

func TestHealthCheck(t *testing.T) {
	RunTest(t, func(ctrl *gomock.Controller) {

		folderChanger := NewFolderChanger("../../cmd/server")
		defer folderChanger.ChangeBack()

		startChan := make(chan bool)

		startup := startup.NewStartup()
		var myEcho *echo.Echo
		hooks := &echo_contracts_startup.Hooks{
			PreStartHook: func(echo *echo.Echo) error {
				myEcho = echo
				startChan <- true
				return nil
			},
		}
		startup.AddHooks(hooks)

		r := runtime.New(startup)
		future := ExecuteWithPromiseAsync(r)

		<-startChan

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		myEcho.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		fmt.Println("data:", string(data))

		r.Stop()
		future.Get()
	})
}
