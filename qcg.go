package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type QcgFacade struct{}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QcgFacade) prepareResource(ids, command string) (string, error) {
	stringOutput, err := execute(ids, command)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	matches := regexp.MustCompile(`jobId = ([\S]+)`).FindStringSubmatch(stringOutput)
	if len(matches) == 0 {
		return "", fmt.Errorf(stringOutput)
	}

	return matches[1], nil
}

//gets resource states
//returns array of resource states
func (qf QcgFacade) StatusCheck() ([]string, error) {
	command := `QCG_ENV_PROXY_DURATION_MIN=12 qcg-list -F "%-25I  %-20S"`

	stringOutput, err := execute("[qcg]", command)
	if err != nil {
		return nil, fmt.Errorf(stringOutput)
	}

	if strings.Contains(stringOutput, "Enter GRID pass phrase for this identity:") {
		logger.Info("Password required, cannot monitor QCG infrastructure\n")
		return nil, fmt.Errorf("Proxy invalid")
	}

	return strings.Split(stringOutput, "\n"), nil
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf QcgFacade) resourceStatus(statusArray []string, jobID string) (string, error) {
	if jobID == "" {
		return "available", nil
	}

	for _, status := range statusArray {
		if strings.Contains(status, jobID) {
			matches := regexp.MustCompile(`(?:\S+\s+)(\S+).+`).FindStringSubmatch(status)
			if len(matches) == 0 {
				return "", fmt.Errorf(status)
			}

			var res string
			switch matches[1] {
			case "UNSUBMITTED":
				{
					res = "initializing"
				}
			case "UNCOMMITED":
				{
					res = "initializing"
				}
			case "QUEUED":
				{
					res = "initializing"
				}
			case "PREPROCESSING":
				{
					res = "initializing"
				}
			case "PENDING":
				{
					res = "initializing"
				}
			case "RUNNING":
				{
					res = "running_sm"
				}
			case "STOPPED":
				{
					res = "released"
				}
			case "POSTPROCESSING":
				{
					res = "released"
				}
			case "FINISHED":
				{
					res = "released"
				}
			case "FAILED":
				{
					res = "released"
				}
			case "CANCELED":
				{
					res = "released"
				}
			case "UNKNOWN":
				{
					res = "error"
				}
			default:
				{
					return "", fmt.Errorf(status)
				}
			}
			return res, nil
		}
	}
	// no such jobID
	return "released", nil
}

//receives sm_record, ExperimentManager connector and infrastructure name
//decides about action on sm and its resources
//returns nothing
func (qf QcgFacade) HandleSM(sm_record *Sm_record, emc *ExperimentManagerConnector, infrastructure string, statusArray []string) {
	resource_status, err := qf.resourceStatus(statusArray, sm_record.Job_id)
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

		//extract first zip
		err := extract("sources_"+sm_record.Id+".zip", ".")
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		//move second zip one directory up
		_, err = executeSilent("mv scalarm_simulation_manager_code_" + sm_record.Sm_uuid + "/* .")
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		//remove both zips and catalog left from first unzip
		_, err = executeSilent("rm -rf  sources_" + sm_record.Id + ".zip scalarm_simulation_manager_code_" + sm_record.Sm_uuid)
		if err != nil {
			sm_record.Error_log = err.Error()
			sm_record.Resource_status = "error"
			return
		}
		logger.Debug("Code files extracted")

		//run command
		jobID, err := qf.prepareResource(ids, sm_record.Cmd_to_execute)
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
