package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

const DEFAULT_PROBE_FREQ_SECS int = 10

func main() {
	// TODO: true versioning (SCAL-937)
	logger.Info("ScalarmWorkersManager 15.06-dev-20150903-1")

	//set config file name
	var configFile string = "config.json"
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	}

	//register working
	RegisterWorking()
	defer UnregisterWorking()

	//declare variables - memory optimization
	var smRecords []SMRecord
	var smRecord SMRecord
	var smRecordOld SMRecord
	var smRecordsRaw interface{}
	var smRecordsCount int

	var infrastructure Infrastructure
	var statusArray []string
	var statusError error
	var err error

	var noMoreRecords bool = false
	var noMoreRecordsTime time.Time

	//listen for signals
	infrastructuresChannel := make(chan []Infrastructure, 10)
	errorChannel := make(chan error, 1)
	go SignalCatcher(infrastructuresChannel, errorChannel, configFile)

	//read configuration
	configData, err := ReadConfiguration(configFile)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Could not read configuration file: %v", configFile))
	}

	//setup verbosity
	logger.SetVerbosity(configData.VerboseMode)

	logger.Info("Config loaded")
	logger.Info("  Information Service address: %v", configData.InformationServiceAddress)
	logger.Info("  Login:                       %v", configData.Login)
	logger.Info("  Scalarm certificate path:    %v", configData.ScalarmCertificatePath)
	logger.Info("  Scalarm scheme:              %v", configData.ScalarmScheme)
	logger.Info("  Insecure SSL:                %v", configData.InsecureSSL)
	logger.Info("  Exit timeout (secs):         %v", configData.ExitTimeoutSecs)
	logger.Info("  Probe frequency (secs):      %v", configData.ProbeFrequencySecs)
	logger.Info("  Verbose mode:                %v", configData.VerboseMode)
	logger.Info("  Infrastructures:             %v", configData.Infrastructures)

	//setup time values
	var waitIndefinitely bool = (configData.ExitTimeoutSecs < 0)
	var exitTimeout time.Duration = time.Duration(configData.ExitTimeoutSecs) * time.Second
	var probeFrequencySecs = time.Duration(DEFAULT_PROBE_FREQ_SECS) * time.Second
	if configData.ProbeFrequencySecs > 0 {
		probeFrequencySecs = time.Duration(configData.ProbeFrequencySecs) * time.Second
	}

	//create EM connector
	emc := NewEMConnector(configData.Login, configData.Password,
		configData.ScalarmCertificatePath, configData.ScalarmScheme, configData.InsecureSSL)

	//get experiment manager location
	if _, err := RepetitiveCaller(
		func() (interface{}, error) {
			return nil, emc.GetExperimentManagerLocation(configData.InformationServiceAddress)
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

		smRecordsCount = 0
		//infrastructures loop
		for _, infrastructure = range configData.Infrastructures {
			//get sm_records
			if smRecordsRaw, err = RepetitiveCaller(
				func() (interface{}, error) {
					return emc.GetSimulationManagerRecords(infrastructure)
				},
				nil,
				"GetSimulationManagerRecords",
			); err != nil {
				logger.Fatal(fmt.Sprintf("Unable to get simulation manager records for %v", infrastructure.Name))
			} else {
				smRecords = smRecordsRaw.([]SMRecord)
			}

			var activeSmRecords []SMRecord

			for _, smRecord := range smRecords {
				if smRecord.State == "error" {
					if smRecord.CmdToExecuteCode != "" {
						activeSmRecords = append(activeSmRecords, smRecord)
					}
				} else {
					activeSmRecords = append(activeSmRecords, smRecord)
				}
			}

			smRecords = activeSmRecords

			logger.Info("[%v] %v sm_records", infrastructure.Name, len(smRecords))
			if len(smRecords) > 0 {
				logger.Debug("\tScalarm ID               Name")
				for _, smRecord = range smRecords {
					logger.Debug("\t%v %v", smRecord.ID, smRecord.Name)
				}
			}

			smRecordsCount += len(smRecords)
			if len(smRecords) == 0 {
				continue
			}

			//check status
			statusArray, statusError = infrastructureFacades[infrastructure.GetInfrastructureName()].StatusCheck()
			if statusError != nil {
				logger.Info("Cannot get status for %v infrastructure", infrastructure.Name)
			}

			//sm_records loop
			for _, smRecord = range smRecords {
				smRecordOld = smRecord

				if statusError != nil {
					//could not read status for current infrastructure
					smRecord.ResourceStatus = "not_available"
					smRecord.ErrorLog = statusError.Error()
				} else {
					//handle SiM

					facade := infrastructureFacades[infrastructure.GetInfrastructureName()]
					err = HandleSiM(facade, &smRecord, infrastructure.GetInfrastructureId(), emc, statusArray)
					if err != nil {
						smRecord.ErrorLog = err.Error()
						smRecord.ResourceStatus = "error"
					}
					// ResourceStatus can be marked to_check after infrastructure action
					if smRecord.ResourceStatus == "to_check" {
						// refresh statusArray
						statusArray, statusError = infrastructureFacades[infrastructure.GetInfrastructureName()].StatusCheck()
						if statusError != nil {
							logger.Info("Cannot get status for %v infrastructure", infrastructure.Name)
						}
						resourceStatus, err := facade.ResourceStatus(statusArray, &smRecord)
						if err != nil {
							smRecord.ErrorLog = err.Error()
							smRecord.ResourceStatus = "error"
						} else {
							smRecord.ResourceStatus = resourceStatus
						}
					}
				}

				//notify state change if needed
				if smRecordOld != smRecord {
					if _, err = RepetitiveCaller(
						func() (interface{}, error) {
							return nil, emc.NotifyStateChange(&smRecord, &smRecordOld, infrastructure.GetInfrastructureId())
						},
						nil,
						"NotifyStateChange",
					); err != nil {
						logger.Fatal("Unable to update simulation manager record")
					}
				}
			}
		}

		//wait for new records if needed
		if !waitIndefinitely && smRecordsCount == 0 {
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
