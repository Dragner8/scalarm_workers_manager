package main

import (
	"fmt"
	"regexp"
	"strings"
)

type SlurmFacade struct {
	Shell 	ShellExecutor
	PLGridFacade
}

//gets resource states
//returns array of resource states
func (qf SlurmFacade) StatusCheck() ([]string, error) {
	command := `squeue -u $USER`

	stringOutput, err := qf.Shell.execute("[slurm]", command)
	if err != nil {
		return nil, fmt.Errorf(stringOutput)
	}

	return strings.Split(stringOutput, "\n"), nil
}

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf SlurmFacade) PrepareResource(ids, command string) (string, error) {
	stringOutput, err := qf.Shell.execute(ids, command)
	if err != nil {
		return "", fmt.Errorf(stringOutput)
	}

	matches := regexp.MustCompile(`(Submitted batch job [\d]+)`).FindStringSubmatch(stringOutput)
	if len(matches) == 0 {
		return "", fmt.Errorf(stringOutput)
	}

	submitJobSentence := matches[1]
	jobId := strings.Fields(submitJobSentence)

	return jobId[len(jobId) - 1], nil
}

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf SlurmFacade) ResourceStatus(statusArray []string, smRecord *SMRecord) (string, error) {
	if smRecord.JobID == "" {
		return "available", nil
	}

	res := ""

	for _, status := range statusArray {
		if strings.Contains(status, smRecord.JobID) {

			for _, statusLineToken := range strings.Fields(status) {
				switch statusLineToken {
					case "PENDING" :
						{
							res = "initializing"
						}
					case "PD":
						{
							res = "initializing"
						}
					case "CONFIGURING":
						{
							res = "initializing"
						}
					case "CF":
						{
							res = "initializing"
						}
					case "RUNNING":
						{
							res = "running_sm"
						}
					case "R":
						{
							res = "running_sm"
						}
					case "COMPLETING":
						{
							res = "released"
						}
					case "CG":
						{
							res = "released"
						}
					case "COMPLETED":
						{
							res = "released"
						}
					case "CD":
						{
							res = "released"
						}
					case "CA":
						{
							res = "released"
						}
					case "CANCELLED":
						{
							res = "released"
						}
					case "FAILED":
						{
							res = "error"
						}
					case "F":
						{
							res = "error"
						}
					case "NODE_FAIL":
						{
							res = "error"
						}
					case "NF":
						{
							res = "error"
						}
					case "PREEMPTED":
						{
							res = "error"
						}
					case "PR":
						{
							res = "error"
						}
					case "SUSPENDED":
						{
							res = "error"
						}
					case "S":
						{
							res = "error"
						}
					case "TIMEOUT":
						{
							res = "error"
						}
					case "TO":
						{
							res = "error"
						}
					default:
						{
							res = ""
						}
				}

				if res != "" {
					return res, nil
				}
			}

			return "", fmt.Errorf(status)
		}
	}

	// no such jobID
	return "released", nil
}
