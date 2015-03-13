package main

import (
	"github.com/scalarm/scalarm_workers_manager/logger"
)

func HandleSiM(facade IInfrastructureFacade, sm_record *Sm_record, emc *EMConnector, infrastructure string, statusArray []string) {
	resource_status, err := facade.ResourceStatus(statusArray, sm_record.Job_id)
	if err != nil {
		sm_record.Error_log = err.Error()
		sm_record.Resource_status = "error"
		return
	}

	ids := sm_record.GetIDs()

	logger.Debug("%v Sm_record state: %v", ids, sm_record.State)
	logger.Debug("%v Resource status: %v", ids, resource_status)

	if sm_record.Cmd_to_execute_code != "" {
		logger.Info("%v Command to execute: %v", ids, sm_record.Cmd_to_execute_code)
	}

	defer func() {
		sm_record.Cmd_to_execute = ""
		sm_record.Cmd_to_execute_code = ""
	}()

	if (sm_record.Cmd_to_execute_code == "prepare_resource" && resource_status == "available") || sm_record.Cmd_to_execute_code == "restart" {

		if _, err := RepetitiveCaller(
			func() (interface{}, error) {
				return nil, emc.GetSimulationManagerCode(sm_record.Id, infrastructure)
			},
			nil,
			"GetSimulationManagerCode",
		); err != nil {
			logger.Fatal("Unable to get simulation manager code")
		}

		facade.ExtractSiMFiles(sm_record)

		//run command
		jobID, err := facade.PrepareResource(ids, sm_record.Cmd_to_execute)
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		logger.Info("%v Assigned job_id: %v", ids, jobID)
		sm_record.Job_id = jobID

	} else if sm_record.Cmd_to_execute_code == "stop" {

		stringOutput, err := execute(ids, sm_record.Cmd_to_execute)
		if err != nil {
			sm_record.Error_log = stringOutput
			sm_record.Resource_status = "error"
			return
		}

	} else if sm_record.Cmd_to_execute_code == "get_log" {

		stringOutput, _ := execute(ids, sm_record.Cmd_to_execute)
		sm_record.Error_log = stringOutput

	}

	sm_record.Resource_status = resource_status
}
