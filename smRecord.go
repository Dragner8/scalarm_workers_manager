package main

import (
	"bytes"

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
	JobID            string `json:"job_id"`
	PID              string `json:"pid"`
	VMID             string `json:"vm_id"`
	ResID            string `json:"res_id"`
}

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
	logger.Debug("\tjob_id              %v", smRecord.JobID)
	logger.Debug("\tpid                 %v", smRecord.PID)
	logger.Debug("\tvm_id               %v", smRecord.VMID)
	logger.Debug("\tres_id              %v", smRecord.ResID)
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
