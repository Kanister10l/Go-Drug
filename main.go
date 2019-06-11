package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kanister10l/Go-Drug/controller"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Println("Error creating logger. Abandoning task!")
		os.Exit(127)
	}
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("App Started", "AppVersion", Version, "BuildTime", Time)

	api := controller.NewApi("", "8080", sugar)

	finish := make(chan bool)
	stopApi := make(chan bool)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)

	api.Listen(finish, stopApi)

	go func() {
		signal := <-signalChannel
		if signal.String() == "interrupt" {
			stopApi <- true
		}
	}()

	<-finish
}
