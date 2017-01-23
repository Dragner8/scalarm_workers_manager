package main

import (
	"reflect"
	"testing"
	// "fmt"
)

func TestSimsGroupingForSingleExperimentByExperimentId(t *testing.T) {
	record1 := SMRecord{ExperimentId: "exp1", ID: "sim1"}
	record2 := SMRecord{ExperimentId: "exp1", ID: "sim2"}
	record3 := SMRecord{ExperimentId: "exp1", ID: "sim3"}

	records := []SMRecord{record1, record2, record3}

	groupingResult := GroupSimsByExperiment(records)

	if len(groupingResult) != 1 {
		t.Errorf("Unexpected value. Got: '%v' - Expected 1", len(groupingResult))
	}

	_, exists := groupingResult["exp1"]

	if !exists {
		t.Errorf("There is no experiment id in the map but there should be!")
	}

	if !reflect.DeepEqual(groupingResult["exp1"], records) {
		t.Errorf("Unexpected value. Got: '%v' - Expected %v", groupingResult["exp1"], records)
	}
}

func TestEmptySimsGroupingByExperimentId(t *testing.T) {
	records := []SMRecord{}

	groupingResult := GroupSimsByExperiment(records)

	if len(groupingResult) != 0 {
		t.Errorf("Unexpected value. Got: '%v' - Expected 1", len(groupingResult))
	}
}

func TestSimsGroupingForManyExperimentsByExperimentId(t *testing.T) {
	record1 := SMRecord{ExperimentId: "exp1", ID: "sim1"}
	record2 := SMRecord{ExperimentId: "exp1", ID: "sim2"}
	record3 := SMRecord{ExperimentId: "exp1", ID: "sim3"}
	record4 := SMRecord{ExperimentId: "exp2", ID: "sim4"}
	record5 := SMRecord{ExperimentId: "exp2", ID: "sim5"}
	record6 := SMRecord{ExperimentId: "exp2", ID: "sim6"}
	record7 := SMRecord{ExperimentId: "exp3", ID: "sim7"}
	record8 := SMRecord{ExperimentId: "exp3", ID: "sim8"}
	record9 := SMRecord{ExperimentId: "exp3", ID: "sim9"}

	records := []SMRecord{record1, record2, record3, record4, record5, record6, record7, record8, record9}

	groupingResult := GroupSimsByExperiment(records)

	if len(groupingResult) != 3 {
		t.Errorf("Unexpected value. Got: '%v' - Expected 1", len(groupingResult))
	}

	_, exists := groupingResult["exp1"]

	if !exists {
		t.Errorf("There is no experiment id in the map but there should be!")
	}

	_, exists = groupingResult["exp2"]

	if !exists {
		t.Errorf("There is no experiment id in the map but there should be!")
	}

	_, exists = groupingResult["exp3"]

	if !exists {
		t.Errorf("There is no experiment id in the map but there should be!")
	}

	if !reflect.DeepEqual(groupingResult["exp1"], []SMRecord{record1, record2, record3}) {
		t.Errorf("Unexpected value. Got: '%v' - Expected %v", groupingResult["exp1"], []SMRecord{record1, record2, record3})
	}

	if !reflect.DeepEqual(groupingResult["exp2"], []SMRecord{record4, record5, record6}) {
		t.Errorf("Unexpected value. Got: '%v' - Expected %v", groupingResult["exp2"], []SMRecord{record4, record5, record6})
	}

	if !reflect.DeepEqual(groupingResult["exp3"], []SMRecord{record7, record8, record9}) {
		t.Errorf("Unexpected value. Got: '%v' - Expected %v", groupingResult["exp3"], []SMRecord{record7, record8, record9})
	}

}

func TestSelectSimsToStart(t *testing.T) {
	record1 := SMRecord{ID: "sim1", CmdToExecuteCode: "restart"}
	record2 := SMRecord{ID: "sim2", CmdToExecuteCode: "get_log#_#restart"}
	record3 := SMRecord{ID: "sim2", CmdToExecuteCode: "prepare_resource"}
	record4 := SMRecord{ID: "sim1", CmdToExecuteCode: "prepare_resource#_#get_log"}
	record5 := SMRecord{ID: "sim2", CmdToExecuteCode: "stop"}
	record6 := SMRecord{ID: "sim2", CmdToExecuteCode: "get_log"}
	record7 := SMRecord{ID: "sim1", CmdToExecuteCode: "destroy"}

	records := []SMRecord{record1, record2, record3, record4, record5, record6, record7}

	filtered_records := SelectStartingSims(records, "available")

	if !reflect.DeepEqual(filtered_records, []SMRecord{record1, record2, record3, record4}) {
		t.Errorf("Unexpected value. Got: '%v' - Expected %v", filtered_records, []SMRecord{record1, record2, record3, record4})
	}

}

func TestSelectSimsToStartWithNotAvailableResource(t *testing.T) {
	record3 := SMRecord{ID: "sim2", CmdToExecuteCode: "prepare_resource"}
	record4 := SMRecord{ID: "sim1", CmdToExecuteCode: "prepare_resource#_#get_log"}
	record5 := SMRecord{ID: "sim2", CmdToExecuteCode: "stop"}
	record6 := SMRecord{ID: "sim2", CmdToExecuteCode: "get_log"}
	record7 := SMRecord{ID: "sim1", CmdToExecuteCode: "destroy"}

	records := []SMRecord{record3, record4, record5, record6, record7}

	filtered_records := SelectStartingSims(records, "not_available")

	if !reflect.DeepEqual(filtered_records, []SMRecord{}) {
		t.Errorf("Unexpected value. Got: '%v' - Expected %v", filtered_records, []SMRecord{})
	}

}

func TestRemoveSimsToStart(t *testing.T) {
	record1 := SMRecord{ID: "sim1", CmdToExecuteCode: "restart"}
	record2 := SMRecord{ID: "sim2", CmdToExecuteCode: "get_log#_#restart"}
	record3 := SMRecord{ID: "sim2", CmdToExecuteCode: "prepare_resource"}
	record4 := SMRecord{ID: "sim1", CmdToExecuteCode: "prepare_resource#_#get_log"}
	record5 := SMRecord{ID: "sim2", CmdToExecuteCode: "stop"}
	record6 := SMRecord{ID: "sim2", CmdToExecuteCode: "get_log"}
	record7 := SMRecord{ID: "sim1", CmdToExecuteCode: "destroy"}

	records := []SMRecord{record1, record2, record3, record4, record5, record6, record7}

	filtered_records := RemoveStartingSims(records, "available")

	if !reflect.DeepEqual(filtered_records, []SMRecord{record5, record6, record7}) {
		t.Errorf("Unexpected value. Got: '%v' - Expected %v", filtered_records, []SMRecord{record5, record6, record7})
	}

}

func TestRemoveSimsToStartWithNotAvailableResource(t *testing.T) {
	record3 := SMRecord{ID: "sim2", CmdToExecuteCode: "prepare_resource"}
	record4 := SMRecord{ID: "sim1", CmdToExecuteCode: "prepare_resource#_#get_log"}
	record5 := SMRecord{ID: "sim2", CmdToExecuteCode: "stop"}
	record6 := SMRecord{ID: "sim2", CmdToExecuteCode: "get_log"}
	record7 := SMRecord{ID: "sim1", CmdToExecuteCode: "destroy"}

	records := []SMRecord{record3, record4, record5, record6, record7}

	filtered_records := RemoveStartingSims(records, "not_available")

	if !reflect.DeepEqual(filtered_records, records) {
		t.Errorf("Unexpected value. Got: '%v' - Expected %v", filtered_records, records)
	}

}
