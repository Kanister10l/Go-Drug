package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"go.uber.org/zap"

	"github.com/go-zoo/bone"
)

type Api struct {
	Sugar *zap.SugaredLogger
	IP    string
	Port  string
}

func NewApi(ip, port string, sugar *zap.SugaredLogger) *Api {
	api := &Api{}
	api.Sugar = sugar
	api.IP = ip
	api.Port = port

	return api
}

func (api *Api) Listen(finish, stop chan bool) {
	go func() {
		defer api.Sugar.Sync()
		mux := bone.New()

		api.Sugar.Info("Api start listening")
		srv := http.Server{
			Addr:    fmt.Sprintf("%s:%s", api.IP, api.Port),
			Handler: mux,
		}

		go func() {
			err := srv.ListenAndServe()
			if err != nil {
				api.Sugar.Fatalw("Api failed to start listening", "Error", err)
				os.Exit(127)
			}
		}()
		<-stop
		srv.Shutdown(context.TODO())
		api.Sugar.Info("Api stopped listening")
		finish <- true
	}()
}
