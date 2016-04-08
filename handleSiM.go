package main

import (
	"fmt"
	"strings"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

const separator = "#_#"

// This function checks if some actions should be invoked based in smRecord information
// It can modify smRecord, in particular: CmdToExecute, CmdToExecuteCode and ResourceStatus
// Notice: This function can set "to_check" ResourceStatus, which is invalid state
// If "to_check" is set, smRecord.ResourceStatus must be set externally
func HandleSiM(facade IInfrastructureFacade, smRecord *SMRecord, infrastructure string, emc *EMConnector, statusArray []string) error {
	var err error

	resourceStatus, err := facade.ResourceStatus(statusArray, smRecord)
	if err != nil {
		return err
	}

	defer func() {
		smRecord.CmdToExecuteCode = ""
		smRecord.CmdToExecute = ""
	}()

	ids := smRecord.GetIDs()

	logger.Debug("%v Sm_record state: %v", ids, smRecord.State)
	logger.Debug("%v Resource status: %v", ids, resourceStatus)

	if smRecord.CmdToExecuteCode == "" {
		// no action was performed, so reourceStatus should be up to date
		smRecord.ResourceStatus = resourceStatus
		return nil
	}

	commandCodes := strings.Split(smRecord.CmdToExecuteCode, separator)
	commands := strings.Split(smRecord.CmdToExecute, separator)

	if len(commands) != len(commandCodes) {
		logger.Info("commands: %v", commands)
		logger.Info("commandCodes: %v", commandCodes)
		return fmt.Errorf("Commands count does not match command codes count")
	}

	for i := 0; i < len(commands); i++ {
		commandCode := commandCodes[i]
		command := commands[i]

		logger.Info("%v Command to execute: %v", ids, commandCode)

		if (commandCode == "prepare_resource" && resourceStatus == "available") || commandCode == "restart" {

			//get code files
			if _, err := RepetitiveCaller(
				func() (interface{}, error) {
					return nil, emc.GetSimulationManagerCode(smRecord.ID, infrastructure)
				},
				nil,
				"GetSimulationManagerCode",
			); err != nil {
				logger.Fatal("Unable to get simulation manager code")
			}

			//extract files
			err = facade.ExtractSiMFiles(smRecord)
			if err != nil {
				return err
			}

			//run command
			jobID, err := facade.PrepareResource(ids, command)
			if err != nil {
				return err
			}
			logger.Info("%v Assigned job_id: %v", ids, jobID)
			facade.SetId(smRecord, jobID)

		} else if commandCode == "stop" {

			stringOutput, err := execute(ids, command)
			if err != nil {
				return fmt.Errorf(stringOutput)
			}

		} else if commandCode == "get_log" {

			stringOutput, _ := execute(ids, command)
			smRecord.ErrorLog = stringOutput

		} else {
			// unknown command
			logger.Info("%v Unsupported command code: \"%v\"", ids, commandCode)
			logger.Info("%v Will NOT execute: \"%v\"", ids, command)
		}

	}

	// an action was performed, so resourceStatus could change
	// set "to_check" flag to indicate need for status check again
	smRecord.ResourceStatus = "to_check"
	return nil
}
