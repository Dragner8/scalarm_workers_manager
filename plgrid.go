package main

import (
	"github.com/scalarm/scalarm_workers_manager/logger"
)

type PLGridFacade struct {
	Name string
}

func (plgf PLGridFacade) ExtractSiMFiles(sm_record *Sm_record) {

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
}
