package main

import "testing"
import "strings"
import "errors"

type SlurmFakeBashExecutor struct {

}

func (SlurmFakeBashExecutor) executeSilent(command string) (string, error) {
  return ``, nil
}

func (SlurmFakeBashExecutor) execute(ids, command string) (string, error) {
  if strings.Contains(command, "squeue -u") {
    return `
blablalba
    JobID    JobName  Partition    Account  AllocCPUS      State ExitCode
------------ ---------- ---------- ---------- ---------- ---------- --------
1131597         TestJob plgrid-te+ scalprome+          1  RUNNING      0:0
`, nil
  }

  if strings.Contains(command, "sbatch") {
    return `
    module load blsabdoasjfasfa
    fFASfdaDDSDA
    Submitted batch job 1131510


    `, nil
  }

  return "", errors.New("")
}

func TestSlurmStatusCheck(t *testing.T) {
  facade := SlurmFacade{SlurmFakeBashExecutor{}, PLGridFacade{}}

  statusArray, _ := facade.StatusCheck()

  expectedList := []string{
    "",
    "blablalba",
    "    JobID    JobName  Partition    Account  AllocCPUS      State ExitCode",
    "------------ ---------- ---------- ---------- ---------- ---------- --------",
    "1131597         TestJob plgrid-te+ scalprome+          1  RUNNING      0:0",
    "",
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

func TestSlurmResourceStatusOfARunningWorker(t *testing.T) {
  facade := SlurmFacade{SlurmFakeBashExecutor{}, PLGridFacade{}}
  statusArray := []string{
    "",
    "blablalba",
    "    JobID    JobName  Partition    Account  AllocCPUS      State ExitCode",
    "------------ ---------- ---------- ---------- ---------- ---------- --------",
    "1131597         TestJob plgrid-te+ scalprome+          1  RUNNING      0:0",
    "",
  }

  smRecord := new(SMRecord)
  smRecord.JobID = "1131597"

  resourceStatus, _ := facade.ResourceStatus(statusArray, smRecord)

  if resourceStatus != "running_sm" {
    t.Errorf("Unexpected value. Got: '%v' - Expected running_sm", resourceStatus)
  }

}

func TestSlurmResourceStatusOfANewWorker(t *testing.T) {
  facade := SlurmFacade{SlurmFakeBashExecutor{}, PLGridFacade{}}
  statusArray := []string{}
  smRecord := new(SMRecord)
  smRecord.JobID = ""

  resourceStatus, _ := facade.ResourceStatus(statusArray, smRecord)

  if resourceStatus != "available" {
    t.Errorf("Unexpected value. Got: '%v' - Expected available", resourceStatus)
  }

}

func TestSlurmPrepareResource(t *testing.T) {
  facade := SlurmFacade{SlurmFakeBashExecutor{}, PLGridFacade{}}
  jobId, _ := facade.PrepareResource("ids", "bash sbatch whatever")

  if jobId != "1131510" {
    t.Errorf("Unexpected value. Got: '%v' - Expected 1131510", jobId)
  }

}
