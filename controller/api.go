package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/kanister10l/Go-Drug/helpers"

	"go.uber.org/zap"

	"github.com/go-zoo/bone"
)

type API struct {
	Sugar *zap.SugaredLogger
	IP    string
	Port  string
}

func NewAPI(ip, port string, sugar *zap.SugaredLogger) *API {
	api := &API{}
	api.Sugar = sugar
	api.IP = ip
	api.Port = port

	return api
}

func (api *API) Listen(ov *helpers.Overwatch) {
	go func() {
		stop := ov.Register()
		defer api.Sugar.Sync()
		mux := bone.New()

		api.Sugar.Info(fmt.Sprintf("API start listening at %s:%s", api.IP, api.Port))
		srv := http.Server{
			Addr:    fmt.Sprintf("%s:%s", api.IP, api.Port),
			Handler: mux,
		}

		go func() {
			err := srv.ListenAndServe()
			if err != nil && err.Error() != "http: Server closed" {
				api.Sugar.Fatalw("API failed to start listening", "Error", err)
				os.Exit(127)
			}
		}()
		<-stop
		srv.Shutdown(context.TODO())
		api.Sugar.Info("API stopped listening")
	}()
}
