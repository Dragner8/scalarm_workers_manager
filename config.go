package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
	//read config file
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	//unmarshal config data
	var configData ConfigData
	err = json.Unmarshal(data, &configData)
	if err != nil {
		return nil, err
	}

	//find and replace "plgrid" infrastructure
	for i, a := range configData.Infrastructures {
		if a == "plgrid" {
			configData.Infrastructures = append(configData.Infrastructures[:i], configData.Infrastructures[i+1:]...)
			configData.Infrastructures, _ = AppendIfMissing(configData.Infrastructures, []string{"qsub", "qcg"})
		}
	}

	//special handling for tilde in path
	if configData.ScalarmCertificatePath != "" {
		if configData.ScalarmCertificatePath[0] == '~' {
			configData.ScalarmCertificatePath = os.Getenv("HOME") + configData.ScalarmCertificatePath[1:]
		}
	}

	//default scheme
	if configData.ScalarmScheme == "" {
		configData.ScalarmScheme = "https"
	}

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
