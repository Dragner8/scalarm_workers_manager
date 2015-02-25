package main

import (
	"log"
	"os"
	"runtime/debug"
	"time"
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
		log.Fatal("Could not read configuration file: " + configFile)
	}

	log.Printf("Config loaded")
	log.Printf("\tInformation Service address: %v", configData.InformationServiceAddress)
	log.Printf("\tLogin:                       %v", configData.Login)
	log.Printf("\tInfrastructures:             %v", configData.Infrastructures)
	log.Printf("\tScalarm certificate path:    %v", configData.ScalarmCertificatePath)
	log.Printf("\tScalarm scheme:              %v", configData.ScalarmScheme)
	log.Printf("\tInsecure SSL:                %v", configData.InsecureSSL)
	log.Printf("\tExit timeout (secs):         %v", configData.ExitTimeoutSecs)
	log.Printf("\tProbe frequency (secs):      %v", configData.ProbeFrequencySecs)
	log.Printf("\tVerbose mode:                %v", configData.VerboseMode)

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
		log.Fatal("Fatal: Unable to get experiment manager location")
	}

	//create infrastructure facades
	infrastructureFacades := NewInfrastructureFacades()

	log.Printf("Configuration finished\n\n\n\n\n")

	for {
		log.Printf("Starting main loop")

		//check for config changes
		configData.Infrastructures = AppendIfMissing(configData.Infrastructures, SignalHandler(infrastructuresChannel, errorChannel))
		log.Printf("Current infrastructures: %v\n\n\n", configData.Infrastructures)

		sm_records_count = 0
		//infrastructures loop
		for _, infrastructure = range configData.Infrastructures {
			log.Printf("Starting " + infrastructure + " infrastructure loop")

			//get sm_records
			if raw_sm_records, err = RepetitiveCaller(
				func() (interface{}, error) {
					return experimentManagerConnector.GetSimulationManagerRecords(infrastructure)
				},
				nil,
				"GetSimulationManagerRecords",
			); err != nil {
				log.Fatal("Fatal: Unable to get simulation manager records for " + infrastructure)
			} else {
				sm_records = raw_sm_records.([]Sm_record)
			}

			sm_records_count += len(sm_records)
			if len(sm_records) == 0 {
				log.Printf("No sm_records")
				break
			}

			//check status
			statusArray, statusError = infrastructureFacades[infrastructure].StatusCheck()

			//sm_records loop
			for _, sm_record = range sm_records {
				old_sm_record = sm_record
				if statusError != nil {
					sm_record.Resource_status = "not_available"
					sm_record.Error_log = statusError.Error()
				} else {
					log.Printf("Starting sm_record handle function, ID: " + sm_record.Id)
					infrastructureFacades[infrastructure].HandleSM(&sm_record, experimentManagerConnector, infrastructure, statusArray)
					log.Printf("Ending sm_record handle function")
				}

				//notify state change
				if old_sm_record != sm_record {
					if _, err = RepetitiveCaller(
						func() (interface{}, error) {
							return nil, experimentManagerConnector.NotifyStateChange(&sm_record, &old_sm_record, infrastructure)
						},
						nil,
						"NotifyStateChange",
					); err != nil {
						log.Fatal("Fatal: Unable to update simulation manager record")
					}
				}
			}
		}

		log.Printf("Ending main loop\n\n\n\n\n")
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
	log.Printf("End")
}
