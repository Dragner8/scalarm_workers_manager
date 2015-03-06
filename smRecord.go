package main

import (
	"bytes"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type Sm_record struct {
	Id                  string `json:"_id"`
	Sm_uuid             string
	State               string
	Resource_status     string
	Cmd_to_execute      string
	Cmd_to_execute_code string
	Error_log           string
	Name                string
	Job_id              string
	Pid                 string
	Vm_id               string
	Res_id              string
}

func (sm Sm_record) Print() {
	logger.Debug("sm_record contents:")
	logger.Debug("\t_id                 %v", sm.Id)
	logger.Debug("\tsm_uuid             %v", sm.Sm_uuid)
	logger.Debug("\tstate               %v", sm.State)
	logger.Debug("\tresource_status     %v", sm.Resource_status)
	logger.Debug("\tcmd_to_execute      %v", sm.Cmd_to_execute)
	logger.Debug("\tcmd_to_execute_code %v", sm.Cmd_to_execute_code)
	logger.Debug("\terror_log           %v", sm.Error_log)
	logger.Debug("\tname                %v", sm.Name)
	logger.Debug("\tjob_id              %v", sm.Job_id)
	logger.Debug("\tpid                 %v", sm.Pid)
	logger.Debug("\tvm_id               %v", sm.Vm_id)
	logger.Debug("\tres_id              %v", sm.Res_id)
}

func (sm Sm_record) GetIDs() string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	buffer.WriteString(sm.Id)
	buffer.WriteString("] [")
	buffer.WriteString(sm.Name)
	buffer.WriteString("]")

	return buffer.String()
}
