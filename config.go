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
	Infrastructures           []Infrastructure
	ScalarmCertificatePath    string
	ScalarmScheme             string
	InsecureSSL               bool
	ExitTimeoutSecs           int
	ProbeFrequencySecs        int
	VerboseMode               bool
}

type Infrastructure struct {
	Name string
	Host string
	Port string
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
		if a.Name == "plgrid" {
			configData.Infrastructures = append(configData.Infrastructures[:i], configData.Infrastructures[i+1:]...)
			configData.Infrastructures, _ = AppendIfMissing(configData.Infrastructures, []Infrastructure{Infrastructure{Name: "qsub"}})
			configData.Infrastructures, _ = AppendIfMissing(configData.Infrastructures, []Infrastructure{Infrastructure{Name: "qcg"}})
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

func innerAppendIfMissing(current []Infrastructure, found Infrastructure) ([]Infrastructure, bool) {
	for _, c := range current {
		if c == found {
			return current, false
		}
	}
	return append(current, found), true
}

func AppendIfMissing(current []Infrastructure, found []Infrastructure) ([]Infrastructure, bool) {
	changed := false
	for _, f := range found {
		change := false
		current, change = innerAppendIfMissing(current, f)
		if change {
			changed = true
		}
	}
	return current, changed
}
