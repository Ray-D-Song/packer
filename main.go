package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"ray-d-song.com/packer/cmd"
	"ray-d-song.com/packer/utils"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			utils.Logger.Error("Recovered from panic", zap.Any("error", r))
			panic(r)
		}
	}()

	go watchExit()
	cmd.Execute()
}

func watchExit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		utils.Logger.Info("Exiting...")
		os.Exit(0)
	}()
}
