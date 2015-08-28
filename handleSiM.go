package main

import (
	"fmt"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

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
	if smRecord.CmdToExecuteCode != "" {
		logger.Info("%v Command to execute: %v", ids, smRecord.CmdToExecuteCode)
	}

	if (smRecord.CmdToExecuteCode == "prepare_resource" && resourceStatus == "available") || smRecord.CmdToExecuteCode == "restart" {

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
		jobID, err := facade.PrepareResource(ids, smRecord.CmdToExecute)
		if err != nil {
			return err
		}
		logger.Info("%v Assigned job_id: %v", ids, jobID)
		facade.SetId(smRecord, jobID)

	} else if smRecord.CmdToExecuteCode == "stop" {

		stringOutput, err := execute(ids, smRecord.CmdToExecute)
		if err != nil {
			return fmt.Errorf(stringOutput)
		}

	} else if smRecord.CmdToExecuteCode == "get_log" {

		stringOutput, _ := execute(ids, smRecord.CmdToExecute)
		smRecord.ErrorLog = stringOutput

	} else {
		// no action was performed, so reourceStatus should not change
		smRecord.ResourceStatus = resourceStatus
		return nil
	}

	// an action was performed, so resourceStatus could change
	// set "to_check" flag to indicate need for status check again 
	smRecord.ResourceStatus = "to_check"
	return nil
}
