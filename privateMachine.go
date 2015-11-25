package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type PrivateMachineFacade struct{}

//gets resource states
//returns array of resource states
func (pmf PrivateMachineFacade) StatusCheck() ([]string, error) {
	command := `ps -A -o pid`

	stringOutput, err := execute("[private machine]", command)
	if err != nil {
		return nil, fmt.Errorf(stringOutput)
	}

	return strings.Split(stringOutput, "\n"), nil
}

//sets id for proper infrastructure
func (pmf PrivateMachineFacade) SetId(smRecord *SMRecord, id string) {
	smRecord.PID = id
}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (pmf PrivateMachineFacade) PrepareResource(ids, command string) (string, error) {
	stringOutput, err := execute(ids, command)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	splitOutput := strings.Split(strings.Trim(stringOutput, " \t\n"), "\n")
	matches := regexp.MustCompile(`\d+`).FindAllString(splitOutput[len(splitOutput)-1], -1)
	if len(matches) == 0 {
		return "", fmt.Errorf(stringOutput)
	}

	return matches[len(matches)-1], nil
}

func (pmf PrivateMachineFacade) ExtractSiMFiles(smRecord *SMRecord) error {

	//extract first zip
	err := extract(fmt.Sprintf("sources_%v.zip", smRecord.ID), ".")
	if err != nil {
		return err
	}

	//move second zip one directory up
	_, err = executeSilent(fmt.Sprintf("mv scalarm_simulation_manager_code_%v/* .", smRecord.SMUUID))
	if err != nil {
		return err
	}

	//remove zip and catalog left from first unzip
	_, err = executeSilent(fmt.Sprintf("rm -rf sources_%v.zip scalarm_simulation_manager_code_%v", smRecord.ID, smRecord.SMUUID))
	if err != nil {
		return err
	}
	logger.Debug("Code files extracted")
	return nil
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (pmf PrivateMachineFacade) ResourceStatus(statusArray []string, smRecord *SMRecord) (string, error) {
	if smRecord.PID == "" {
		return "available", nil
	}

	for _, status := range statusArray {
		status = strings.Trim(status, " \t\n")
		if smRecord.PID == status {
			return "running_sm", nil
		}
	}
	// no such jobID
	return "released", nil
}
