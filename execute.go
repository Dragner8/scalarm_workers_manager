package main

import (
	"os/exec"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

func executeSilent(command string) (string, error) {
	output, err := exec.Command("bash", "-c", command).CombinedOutput()
	stringOutput := string(output[:])

	return stringOutput, err
}

func execute(ids, command string) (string, error) {
	logger.Debug("%v Executing: %v", ids, command)
	stringOutput, err := executeSilent(command)
	logger.Debug("%v Response: %v", ids, stringOutput)
	return stringOutput, err
}

// func scriptExecute(command string) (string, error) {

// 	cmd := []byte("#!/bin/bash\n" + command + "\n")
// 	ioutil.WriteFile("./s.sh", cmd, 0755)
// 	output, err := exec.Command("./s.sh").CombinedOutput()
// 	stringOutput := string(output[:])
// 	os.Remove("s.sh")

// 	return stringOutput, err
// }
