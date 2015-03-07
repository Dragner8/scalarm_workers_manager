package main

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

const DEFAULT_PROBE_FREQ_SECS int = 10

func main() {

	//set config file name
	var configFile string = "config.json"
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	}

	//register working
	RegisterWorking()
	defer UnregisterWorking()

	//declare variables - memory optimization
	var sm_records []Sm_record
	var sm_record Sm_record
	var old_sm_record Sm_record
	var raw_sm_records interface{}
	var sm_records_count int

	var infrastructure string
	var statusArray []string
	var statusError error
	var err error

	var noMoreRecords bool = false
	var noMoreRecordsTime time.Time

	//listen for signals
	infrastructuresChannel := make(chan []string, 10)
	errorChannel := make(chan error, 1)
	go SignalCatcher(infrastructuresChannel, errorChannel, configFile)

	//read configuration
	configData, err := ReadConfiguration(configFile)
	if err != nil {
		logger.Fatal("Could not read configuration file: " + configFile)
	}

	//setup verbosity
	logger.SetVerbosity(configData.VerboseMode)

	logger.Info("Config loaded")
	logger.Info("\tInformation Service address: %v", configData.InformationServiceAddress)
	logger.Info("\tLogin:                       %v", configData.Login)
	logger.Info("\tInfrastructures:             %v", configData.Infrastructures)
	logger.Info("\tScalarm certificate path:    %v", configData.ScalarmCertificatePath)
	logger.Info("\tScalarm scheme:              %v", configData.ScalarmScheme)
	logger.Info("\tInsecure SSL:                %v", configData.InsecureSSL)
	logger.Info("\tExit timeout (secs):         %v", configData.ExitTimeoutSecs)
	logger.Info("\tProbe frequency (secs):      %v", configData.ProbeFrequencySecs)
	logger.Info("\tVerbose mode:                %v", configData.VerboseMode)

	//setup time values
	var waitIndefinitely bool = (configData.ExitTimeoutSecs < 0)
	var exitTimeout time.Duration = time.Duration(configData.ExitTimeoutSecs) * time.Second
	var probeFrequencySecs = time.Duration(DEFAULT_PROBE_FREQ_SECS) * time.Second
	if configData.ProbeFrequencySecs > 0 {
		probeFrequencySecs = time.Duration(configData.ProbeFrequencySecs) * time.Second
	}

	//create EM connector
	experimentManagerConnector := NewExperimentManagerConnector(configData.Login, configData.Password,
		configData.ScalarmCertificatePath, configData.ScalarmScheme, configData.InsecureSSL)

	//get experiment manager location
	if _, err := RepetitiveCaller(
		func() (interface{}, error) {
			return nil, experimentManagerConnector.GetExperimentManagerLocation(configData.InformationServiceAddress)
		},
		nil,
		"GetExperimentManagerLocation",
	); err != nil {
		logger.Fatal("Unable to get experiment manager location")
	}

	//create infrastructure facades
	infrastructureFacades := NewInfrastructureFacades()

	logger.Info("Configuration finished")

	for {
		//check for config changes
		changed := false
		configData.Infrastructures, changed = AppendIfMissing(configData.Infrastructures, SignalHandler(infrastructuresChannel, errorChannel))
		if changed {
			logger.Info("Infrastractures reloaded, current infrastructures: %v", configData.Infrastructures)
		}

		sm_records_count = 0
		//infrastructures loop
		for _, infrastructure = range configData.Infrastructures {
			//get sm_records
			if raw_sm_records, err = RepetitiveCaller(
				func() (interface{}, error) {
					return experimentManagerConnector.GetSimulationManagerRecords(infrastructure)
				},
				nil,
				"GetSimulationManagerRecords",
			); err != nil {
				logger.Fatal("Unable to get simulation manager records for " + infrastructure)
			} else {
				sm_records = raw_sm_records.([]Sm_record)
			}

			logger.Info("[%v] %v sm_records", infrastructure, len(sm_records))
			if len(sm_records) > 0 {
				logger.Debug("\tScalarm ID               Name")
				for _, sm_record = range sm_records {
					logger.Debug("\t%v %v", sm_record.Id, sm_record.Name)
				}
			}

			sm_records_count += len(sm_records)
			if len(sm_records) == 0 {
				continue
			}

			//check status
			statusArray, statusError = infrastructureFacades[infrastructure].StatusCheck()
			if statusError != nil {
				logger.Info("Cannot get status for %v infrastructure", infrastructure)
			}

			//sm_records loop
			for _, sm_record = range sm_records {
				old_sm_record = sm_record
				if statusError != nil {
					sm_record.Resource_status = "not_available"
					sm_record.Error_log = statusError.Error()
				} else {
					infrastructureFacades[infrastructure].HandleSM(&sm_record, experimentManagerConnector, infrastructure, statusArray)
				}

				//notify state change
				if old_sm_record != sm_record {
					_, err = RepetitiveCaller(
						func() (interface{}, error) {
							return nil, experimentManagerConnector.NotifyStateChange(&sm_record, &old_sm_record, infrastructure)
						},
						nil,
						"NotifyStateChange",
					)
					if err != nil {
						logger.Fatal("Unable to update simulation manager record")
					}
				}
			}
		}

		if !waitIndefinitely && sm_records_count == 0 {
			if !noMoreRecords {
				noMoreRecords = true
				noMoreRecordsTime = time.Now()
			}

			if time.Now().After(noMoreRecordsTime.Add(exitTimeout)) {
				break
			}
		} else {
			noMoreRecords = false
		}

		debug.FreeOSMemory()
		time.Sleep(probeFrequencySecs)
	}
	logger.Info("End")
}
