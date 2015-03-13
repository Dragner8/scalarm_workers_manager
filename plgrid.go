package main

import (
	"github.com/scalarm/scalarm_workers_manager/logger"
)

type PLGridFacade struct {
	Name string
}

func (plgf PLGridFacade) ExtractSiMFiles(smRecord *SMRecord) {

	//extract first zip
	err := extract("sources_"+smRecord.ID+".zip", ".")
	if err != nil {
		smRecord.ErrorLog = err.Error()
		smRecord.ResourceStatus = "error"
		return
	}
	//move second zip one directory up
	_, err = executeSilent("mv scalarm_simulation_manager_code_" + smRecord.SMUUID + "/* .")
	if err != nil {
		smRecord.ErrorLog = err.Error()
		smRecord.ResourceStatus = "error"
		return
	}
	//remove both zips and catalog left from first unzip
	_, err = executeSilent("rm -rf  sources_" + smRecord.ID + ".zip scalarm_simulation_manager_code_" + smRecord.SMUUID)
	if err != nil {
		smRecord.ErrorLog = err.Error()
		smRecord.ResourceStatus = "error"
		return
	}
	logger.Debug("Code files extracted")
}
