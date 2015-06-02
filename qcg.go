package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type QcgFacade struct {
	PLGridFacade
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

//receives command to execute
//executes command, extracts resource ID
//returns new job ID
func (qf QcgFacade) PrepareResource(ids, command string) (string, error) {
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

//receives job ID
//checks resource state based on job state
//returns resource state
func (qf QcgFacade) ResourceStatus(statusArray []string, smRecord *SMRecord) (string, error) {
	if smRecord.JobID == "" {
		return "available", nil
	}

	for _, status := range statusArray {
		if strings.Contains(status, smRecord.JobID) {
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
