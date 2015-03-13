package main

import (
	"github.com/scalarm/scalarm_workers_manager/logger"
)

func HandleSiM(facade IInfrastructureFacade, smRecord *SMRecord, emc *EMConnector, infrastructure string, statusArray []string) {
	resourceStatus, err := facade.ResourceStatus(statusArray, smRecord)
	if err != nil {
		smRecord.ErrorLog = err.Error()
		smRecord.ResourceStatus = "error"
		return
	}

	ids := smRecord.GetIDs()

	logger.Debug("%v Sm_record state: %v", ids, smRecord.State)
	logger.Debug("%v Resource status: %v", ids, resourceStatus)

	if smRecord.CmdToExecuteCode != "" {
		logger.Info("%v Command to execute: %v", ids, smRecord.CmdToExecuteCode)
	}

	defer func() {
		smRecord.CmdToExecuteCode = ""
		smRecord.CmdToExecute = ""
	}()

	if (smRecord.CmdToExecuteCode == "prepare_resource" && resourceStatus == "available") || smRecord.CmdToExecuteCode == "restart" {

		if _, err := RepetitiveCaller(
			func() (interface{}, error) {
				return nil, emc.GetSimulationManagerCode(smRecord.ID, infrastructure)
			},
			nil,
			"GetSimulationManagerCode",
		); err != nil {
			logger.Fatal("Unable to get simulation manager code")
		}

		facade.ExtractSiMFiles(smRecord)

		//run command
		jobID, err := facade.PrepareResource(ids, smRecord.CmdToExecute)
		if err != nil {
			smRecord.ErrorLog = err.Error()
			smRecord.ResourceStatus = "error"
			return
		}
		logger.Info("%v Assigned job_id: %v", ids, jobID)
		smRecord.JobID = jobID

	} else if smRecord.CmdToExecuteCode == "stop" {

		stringOutput, err := execute(ids, smRecord.CmdToExecute)
		if err != nil {
			smRecord.ErrorLog = stringOutput
			smRecord.ResourceStatus = "error"
			return
		}

	} else if smRecord.CmdToExecuteCode == "get_log" {

		stringOutput, _ := execute(ids, smRecord.CmdToExecute)
		smRecord.ErrorLog = stringOutput

	}

	smRecord.ResourceStatus = resourceStatus
}
