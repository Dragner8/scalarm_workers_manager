package main

import (
	"fmt"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type PLGridFacade struct{}

//sets id for proper infrastructure
func (plgf PLGridFacade) SetId(smRecord *SMRecord, id string) {
	smRecord.JobID = id
}

func (plgf PLGridFacade) ExtractSiMFiles(smRecord *SMRecord) error {

	//extract first zip
	err := extract(fmt.Sprintf("sources_%v.zip", smRecord.ID), ".")
	if err != nil {
		return err
	}
	//move catalog contents one directory up
	_, err = executeSilent(fmt.Sprintf("mv scalarm_simulation_manager_code_%v/* .", smRecord.SMUUID))
	if err != nil {
		return err
	}
	//remove first zip and catalog left from first unzip
	_, err = executeSilent(fmt.Sprintf("rm -rf  sources_%v.zip scalarm_simulation_manager_code_%v", smRecord.ID, smRecord.SMUUID))
	if err != nil {
		return err
	}
	logger.Debug("Code files extracted")
	return nil
}
