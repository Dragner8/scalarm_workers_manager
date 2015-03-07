package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

func RegisterWorking() {
	logger.Info("Checking for working monitoring mark")
	if _, err := os.Stat(".monitoring_working_mark"); err == nil {
		logger.Info("Mark file exists...")
		pid, _ := ioutil.ReadFile(".monitoring_working_mark")
		output, _ := exec.Command("bash", "-c", "ps -p "+string(pid[:])+" | tail -n +2").CombinedOutput()
		if strings.Contains(string(output[:]), "scalarm") {
			logger.Info("...and process with saved pid [%s] is working:\n%v", string(pid[:]), string(output[:]))
			exec.Command("bash", "-c", "kill -USR1 "+string(pid[:])).Run()
			logger.Fatal("Monitoring already working")
		}
		logger.Info("...but no process with saved pid [%s] is working", string(pid[:]))
	}

	pid := []byte(strconv.Itoa(os.Getpid()))
	logger.Info("Creating monitoring mark file, pid: %s", pid)
	ioutil.WriteFile(".monitoring_working_mark", pid, 0644)
}

func UnregisterWorking() {
	logger.Info("Deleting monitoring mark file")
	os.Remove(".monitoring_working_mark")
}
