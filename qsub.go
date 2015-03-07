package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type QsubFacade struct{}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QsubFacade) prepareResource(ids, command string) (string, error) {
	stringOutput, err := execute(ids, command)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	matches := regexp.MustCompile(`([\d]+.batch.grid.cyf-kr.edu.pl)`).FindStringSubmatch(stringOutput)
	if len(matches) == 0 {
		return "", fmt.Errorf(stringOutput)
	}

	return matches[1], nil
}

//gets resource states
//returns array of resource states
func (qf QsubFacade) StatusCheck() ([]string, error) {
	command := `qstat -u $USER`

	stringOutput, err := execute("[qsub]", command)
	if err != nil {
		return nil, fmt.Errorf(stringOutput)
	}

	return strings.Split(stringOutput, "\n"), nil
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf QsubFacade) resourceStatus(statusArray []string, jobID string) (string, error) {
	if jobID == "" {
		return "available", nil
	}

	for _, status := range statusArray {
		if strings.Contains(status, strings.Split(jobID, ".")[0]) {
			matches := regexp.MustCompile(`(?:\S+\s+){9}([A-Z]).+`).FindStringSubmatch(status)
			if len(matches) == 0 {
				return "", fmt.Errorf(status)
			}

			var res string
			switch matches[1] {
			case "Q":
				{
					res = "initializing"
				}
			case "W":
				{
					res = "initializing"
				}
			case "H":
				{
					res = "running_sm"
				}
			case "R":
				{
					res = "running_sm"
				}
			case "T":
				{
					res = "running_sm"
				}
			case "C":
				{
					res = "released"
				}
			case "E":
				{
					res = "released"
				}
			case "U":
				{
					res = "released"
				}
			case "S":
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
func (qf QsubFacade) HandleSM(sm_record *Sm_record, emc *ExperimentManagerConnector, infrastructure string, statusArray []string) {
	resource_status, err := qf.resourceStatus(statusArray, sm_record.Job_id)
	if err != nil {
		sm_record.Error_log = err.Error()
		sm_record.Resource_status = "error"
		return
	}

	ids := sm_record.GetIDs()

	logger.Debug("%v Sm_record state: %v ", ids, sm_record.State)
	logger.Debug("%v Resource status: %v ", ids, resource_status)

	if sm_record.Cmd_to_execute_code != "" {
		logger.Info("%v Command to execute: %v ", ids, sm_record.Cmd_to_execute_code)
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
		logger.Debug("%v Code files extracted", ids)

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
