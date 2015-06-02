package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

func SignalCatcher(infrastructuresChannel chan []Infrastructure, errorChannel chan error, configFile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)

	for {
		<-c
		newConfig, err := ReadConfiguration(configFile)
		if err != nil {
			errorChannel <- err
		}
		infrastructuresChannel <- newConfig.Infrastructures
	}
}

func SignalHandler(infrastructuresChannel chan []Infrastructure, errorChannel chan error) []Infrastructure {
	//check for errors
	select {
	case err, ok := <-errorChannel:
		if ok {
			logger.Info("An error occured while reloading config: ", err.Error())
		} else {
			logger.Fatal("Channel closed!")
		}
	default:
	}

	//check for config changes
	select {
	case addedInfrastructures, ok := <-infrastructuresChannel:
		if ok {
			return addedInfrastructures
		} else {
			logger.Fatal("Channel closed!")
		}
	default:
	}

	return nil
}
