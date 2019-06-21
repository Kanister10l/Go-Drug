package main

import (
	"github.com/kanister10l/Go-Drug/helpers"
)

//go:generate scripts/build_version.sh

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanister10l/Go-Drug/controller"
	"go.uber.org/zap"
)

func main() {
	sugar := setupLogger()
	sp := helpers.NewParams()
	ov := helpers.NewOverwatch(sugar, 10)

	setupAPI(ov, sp, sugar)
	setupSignal(ov, sugar)

	sugar.Infow("App Started", "AppVersion", Version, "BuildTime", Time)
	<-ov.Final
}

func setupAPI(ov *helpers.Overwatch, sp *helpers.StartupParams, sugar *zap.SugaredLogger) {
	api := controller.NewAPI(*sp.IP, *sp.Port, sugar)
	api.Listen(ov)
}

func setupLogger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Println("Error creating logger. Abandoning task!")
		os.Exit(127)
	}
	defer logger.Sync()

	return logger.Sugar()
}

func setupSignal(ov *helpers.Overwatch, sugar *zap.SugaredLogger) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)

	go func() {
		signal := <-signalChannel
		sugar.Info(fmt.Sprintf("Captured %s Signal", signal.String()))
		if signal.String() == "interrupt" {
			ov.SigInt <- true
		}
	}()
}
