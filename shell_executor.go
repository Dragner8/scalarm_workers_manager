package main

import (
	"os/exec"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type ShellExecutor interface {
  execute(ids, command string) (string, error)
  executeSilent(command string) (string, error)
}

type BashExecutor struct {

}

func (BashExecutor) executeSilent(command string) (string, error) {
	output, err := exec.Command("bash", "-c", command).CombinedOutput()
	stringOutput := string(output[:])

	return stringOutput, err
}

func (BashExecutor) execute(ids, command string) (string, error) {
	logger.Info("%v Executing: %v", ids, command)
	stringOutput, err := executeSilent(command)
	logger.Debug("%v Response: %v", ids, stringOutput)
	return stringOutput, err
}
