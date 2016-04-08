package main

import "testing"
import "fmt"

type FakeBashExecutor struct {

}

func (FakeBashExecutor) executeSilent(command string) (string, error) {
  return `
batch.grid.cyf-kr.edu.pl:
                                                                         Req'd  Req'd   Elap
Job ID               Username    Queue    Jobname          SessID NDS   TSK    Memory Time  S Time
-------------------- ----------- -------- ---------------- ------ ----- ------ ------ ----- - -----
65174528.batch.g     plgkrol     plgrid   STDIN             28652   --     --     --  72:00 R 69:44`, nil
}

func (FakeBashExecutor) execute(ids, command string) (string, error) {
  return `
batch.grid.cyf-kr.edu.pl:
                                                                         Req'd  Req'd   Elap
Job ID               Username    Queue    Jobname          SessID NDS   TSK    Memory Time  S Time
-------------------- ----------- -------- ---------------- ------ ----- ------ ------ ----- - -----
65174528.batch.g     plgkrol     plgrid   STDIN             28652   --     --     --  72:00 R 69:44`, nil
}

func TestQsubStatusCheck(t *testing.T) {
  fmt.Println("Running TestQsubStatusCheck")
  facade := QsubFacade{FakeBashExecutor{}, PLGridFacade{}}

  statusArray, _ := facade.StatusCheck()

  expectedList := []string{
    "",
    "batch.grid.cyf-kr.edu.pl:",
    "                                                                         Req'd  Req'd   Elap",
    "Job ID               Username    Queue    Jobname          SessID NDS   TSK    Memory Time  S Time",
    "-------------------- ----------- -------- ---------------- ------ ----- ------ ------ ----- - -----",
    "65174528.batch.g     plgkrol     plgrid   STDIN             28652   --     --     --  72:00 R 69:44",
  }

  if len(statusArray) != len(expectedList) {
		t.Errorf("Unexpected length. Got: '%d' - Expected %d", len(statusArray), len(expectedList))
	}

  for i := 0; i < len(expectedList); i++ {
    if statusArray[i] != expectedList[i] {
      t.Errorf("Unexpected line Got: '%s' - Expected %s", statusArray[i], expectedList[i])
    }
  }
}

func TestResourceStatusOfARunningWorker(t *testing.T) {
  fmt.Println("Running TestResourceStatusOfARunningWorker")
  facade := QsubFacade{FakeBashExecutor{}, PLGridFacade{}}
  statusArray := []string{
    "",
    "batch.grid.cyf-kr.edu.pl:",
    "                                                                         Req'd  Req'd   Elap",
    "Job ID               Username    Queue    Jobname          SessID NDS   TSK    Memory Time  S Time",
    "-------------------- ----------- -------- ---------------- ------ ----- ------ ------ ----- - -----",
    "65174528.batch.g     plgkrol     plgrid   STDIN             28652   --     --     --  72:00 R 69:44",
  }

  smRecord := new(SMRecord)
  smRecord.JobID = "65174528"

  resourceStatus, _ := facade.ResourceStatus(statusArray, smRecord)

  if resourceStatus != "running_sm" {
    t.Errorf("Unexpected value. Got: '%v' - Expected running_sm", resourceStatus)
  }

}

func TestResourceStatusOfANewWorker(t *testing.T) {
  fmt.Println("Running TestResourceStatusOfANewWorker")
  facade := QsubFacade{FakeBashExecutor{}, PLGridFacade{}}
  statusArray := []string{}
  smRecord := new(SMRecord)
  smRecord.JobID = ""

  resourceStatus, _ := facade.ResourceStatus(statusArray, smRecord)

  if resourceStatus != "available" {
    t.Errorf("Unexpected value. Got: '%v' - Expected available", resourceStatus)
  }

}
