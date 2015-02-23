package main

import "time"

func Heartbeat(experimentManagerConnector *ExperimentManagerConnector, infrastructures []string, hbchan chan []string) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			select {
			case infrastructures = <-hbchan:
			default:
			}
			experimentManagerConnector.Ping(infrastructures)
		}
	}
}
