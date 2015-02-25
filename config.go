package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type ConfigData struct {
	InformationServiceAddress string
	Login                     string
	Password                  string
	Infrastructures           []string
	ScalarmCertificatePath    string
	ScalarmScheme             string
	InsecureSSL               bool
	ExitTimeoutSecs           int
	ProbeFrequencySecs        int
	VerboseMode               bool
}

func ReadConfiguration(configFile string) (*ConfigData, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var configData ConfigData
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return nil, err
	}

	for i, a := range configData.Infrastructures {
		if a == "plgrid" {
			configData.Infrastructures = append(configData.Infrastructures[:i], configData.Infrastructures[i+1:]...)
			configData.Infrastructures, _ = AppendIfMissing(configData.Infrastructures, []string{"qsub", "qcg"})
		}
	}

	if configData.ScalarmCertificatePath != "" {
		if configData.ScalarmCertificatePath[0] == '~' {
			configData.ScalarmCertificatePath = os.Getenv("HOME") + configData.ScalarmCertificatePath[1:]
		}
	}

	if configData.ScalarmScheme == "" {
		configData.ScalarmScheme = "https"
	}

	VERBOSE = configData.VerboseMode

	return &configData, nil
}

func innerAppendIfMissing(currentInfrastructures []string, newInfrastructure string) ([]string, bool) {
	for _, c := range currentInfrastructures {
		if c == newInfrastructure {
			return currentInfrastructures, false
		}
	}
	return append(currentInfrastructures, newInfrastructure), true
}

func AppendIfMissing(currentInfrastructures []string, newInfrastructures []string) ([]string, bool) {
	changed := false
	for _, n := range newInfrastructures {
		change := false
		currentInfrastructures, change = innerAppendIfMissing(currentInfrastructures, n)
		if change {
			changed = true
		}
	}
	return currentInfrastructures, changed
}

func SignalCatcher(infrastructuresChannel chan []string, errorChannel chan error, configFile string) {
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

func SignalHandler(infrastructuresChannel chan []string, errorChannel chan error) []string {
	//check for errors
	select {
	case err, ok := <-errorChannel:
		if ok {
			log.Printf("An error occured while reloading config: " + err.Error())
		} else {
			log.Fatal("Channel closed!")
		}
	default:
	}

	//check for config changes
	select {
	case addedInfrastructures, ok := <-infrastructuresChannel:
		if ok {
			return addedInfrastructures
		} else {
			log.Fatal("Channel closed!")
		}
	default:
	}

	return nil
}
