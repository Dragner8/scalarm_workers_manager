package main

import (
	"bytes"
	"strings"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type SMRecord struct {
	ID               string `json:"_id"`
	SMUUID           string `json:"sm_uuid"`
	State            string `json:"state"`
	ResourceStatus   string `json:"resource_status"`
	CmdToExecute     string `json:"cmd_to_execute"`
	CmdToExecuteCode string `json:"cmd_to_execute_code"`
	ErrorLog         string `json:"error_log"`
	Name             string `json:"name"`
	JobID            string `json:"job_identifier"`
	PID              string `json:"pid"`
	VMID             string `json:"vm_identifier"`
	ResID            string `json:"res_name"`
	ExperimentId     string `json:"experiment_id"`
}

const cmd_separator = "#_#"

func (smRecord SMRecord) Print() {
	logger.Debug("sm_record contents:")
	logger.Debug("\t_id                 %v", smRecord.ID)
	logger.Debug("\tsm_uuid             %v", smRecord.SMUUID)
	logger.Debug("\tstate               %v", smRecord.State)
	logger.Debug("\tresource_status     %v", smRecord.ResourceStatus)
	logger.Debug("\tcmd_to_execute      %v", smRecord.CmdToExecute)
	logger.Debug("\tcmd_to_execute_code %v", smRecord.CmdToExecuteCode)
	logger.Debug("\terror_log           %v", smRecord.ErrorLog)
	logger.Debug("\tname                %v", smRecord.Name)
	logger.Debug("\tjob_identifier      %v", smRecord.JobID)
	logger.Debug("\tpid                 %v", smRecord.PID)
	logger.Debug("\tvm_identifier       %v", smRecord.VMID)
	logger.Debug("\tres_name            %v", smRecord.ResID)
}

func (smRecord SMRecord) GetIDs() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	buffer.WriteString(smRecord.ID)
	buffer.WriteString("] [")
	buffer.WriteString(smRecord.Name)
	buffer.WriteString("]")

	return buffer.String()
}

func (smRecord SMRecord) IsAboutToStart(resourceStatus string) bool {
	commandCodes := strings.Split(smRecord.CmdToExecuteCode, cmd_separator)

	for _, commandCode := range commandCodes {
		if (commandCode == "prepare_resource" && resourceStatus == "available") || commandCode == "restart" {
			return true
		}
	}

	return false
}

// functions operating on data structures with SMRecord
func GroupSimsByExperiment(records []SMRecord) map[string][]SMRecord {
	groupedSims := make(map[string][]SMRecord)

	for _, record := range records {
		_, exists := groupedSims[record.ExperimentId]

		if !exists {
			groupedSims[record.ExperimentId] = make([]SMRecord, 0)
		}

		groupedSims[record.ExperimentId] = append(groupedSims[record.ExperimentId], record)
	}

	return groupedSims
}

func FilterStartingSims(records []SMRecord, resourceStatus string) []SMRecord {
	filteredRecords := make([]SMRecord, 0)

	for _, record := range records {
		if record.IsAboutToStart(resourceStatus) {
			filteredRecords = append(filteredRecords, record)
		}
	}

	return filteredRecords
}
