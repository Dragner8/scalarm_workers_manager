package main

import (
	"fmt"
	"regexp"
	"strings"
)

type QsubFacade struct {
	PLGridFacade
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

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QsubFacade) PrepareResource(ids, command string) (string, error) {
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

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf QsubFacade) ResourceStatus(statusArray []string, smRecord *SMRecord) (string, error) {
	if smRecord.JobID == "" {
		return "available", nil
	}

	for _, status := range statusArray {
		if strings.Contains(status, strings.Split(smRecord.JobID, ".")[0]) {
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
