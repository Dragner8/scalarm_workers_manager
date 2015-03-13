package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"

	"github.com/scalarm/scalarm_workers_manager/logger"
)

type EMConnector struct {
	login                    string
	password                 string
	experimentManagerAddress string
	client                   *http.Client
	scheme                   string
}

func NewEMConnector(login, password, certificatePath, scheme string, insecure bool) *EMConnector {
	var client *http.Client
	tlsConfig := tls.Config{InsecureSkipVerify: insecure}

	if certificatePath != "" {
		CA_Pool := x509.NewCertPool()
		severCert, err := ioutil.ReadFile(certificatePath)
		if err != nil {
			logger.Fatal("Could not load Scalarm certificate")
		}
		CA_Pool.AppendCertsFromPEM(severCert)

		tlsConfig.RootCAs = CA_Pool
	}

	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tlsConfig}}

	return &EMConnector{login: login, password: password, client: client, scheme: scheme}
}

func (emc *EMConnector) GetExperimentManagerLocation(informationServiceAddress string) error {
	resp, err := emc.client.Get(fmt.Sprintf("%v://%v/experiment_managers", emc.scheme, informationServiceAddress))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var experimentManagerAddresses []string
	err = json.Unmarshal(body, &experimentManagerAddresses)
	if err != nil {
		return err
	}

	emc.experimentManagerAddress = experimentManagerAddresses[0] // TODO random

	return nil
}

type EMJsonResponse struct {
	Status    string     `json:"status"`
	SMRecords []SMRecord `json:"sm_records"`
}

func (emc *EMConnector) GetSimulationManagerRecords(infrastructure string) ([]SMRecord, error) {
	urlString := fmt.Sprintf("%v://%v/simulation_managers?", emc.scheme, emc.experimentManagerAddress)
	params := url.Values{}
	params.Add("infrastructure", infrastructure)
	params.Add("options", "{\"states_not\":\"error\",\"onsite_monitoring\":true}")
	urlString = urlString + params.Encode()

	request, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(emc.login, emc.password)

	resp, err := emc.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response EMJsonResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Status != "ok" {
		return nil, errors.New("Damaged data")
	}

	return response.SMRecords, nil
}

func (emc *EMConnector) GetSimulationManagerCode(smRecordId string, infrastructure string) error {
	debug.FreeOSMemory()
	urlString := fmt.Sprintf("%v://%v/simulation_managers/%v/code?", emc.scheme, emc.experimentManagerAddress, smRecordId)
	params := url.Values{}
	params.Add("infrastructure", infrastructure)
	urlString = urlString + params.Encode()

	request, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return err
	}
	request.SetBasicAuth(emc.login, emc.password)

	resp, err := emc.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("sources_%v.zip", smRecordId), body, 0600)
	if err != nil {
		return err
	}

	return nil
}

func inner_sm_record_marshal(current, old, name string, comma *bool, parameters *bytes.Buffer) {
	if current != old {
		if *comma {
			parameters.WriteString(",")
		}
		parameters.WriteString("\"" + name + "\":\"" + escape(current) + "\"")
		*comma = true
	}
}

func sm_record_marshal(smRecord, smRecordOld *SMRecord) string {
	var parameters bytes.Buffer
	parameters.WriteString("{")
	comma := false

	inner_sm_record_marshal(smRecord.SMUUID, smRecordOld.SMUUID, "sm_uuid", &comma, &parameters)

	inner_sm_record_marshal(smRecord.State, smRecordOld.State, "state", &comma, &parameters)

	inner_sm_record_marshal(smRecord.ResourceStatus, smRecordOld.ResourceStatus, "resource_status", &comma, &parameters)

	inner_sm_record_marshal(smRecord.CmdToExecute, smRecordOld.CmdToExecute, "cmd_to_execute", &comma, &parameters)

	inner_sm_record_marshal(smRecord.CmdToExecuteCode, smRecordOld.CmdToExecuteCode, "cmd_to_execute_code", &comma, &parameters)

	inner_sm_record_marshal(smRecord.ErrorLog, smRecordOld.ErrorLog, "error_log", &comma, &parameters)

	inner_sm_record_marshal(smRecord.Name, smRecordOld.Name, "name", &comma, &parameters)

	inner_sm_record_marshal(smRecord.JobID, smRecordOld.JobID, "job_id", &comma, &parameters)

	inner_sm_record_marshal(smRecord.PID, smRecordOld.PID, "pid", &comma, &parameters)

	inner_sm_record_marshal(smRecord.VMID, smRecordOld.VMID, "vm_id", &comma, &parameters)

	inner_sm_record_marshal(smRecord.ResID, smRecordOld.ResID, "res_id", &comma, &parameters)

	parameters.WriteString("}")

	logger.Info("%v Update: %v", smRecord.GetIDs(), parameters.String())
	return parameters.String()
}

func (emc *EMConnector) NotifyStateChange(smRecord, smRecordOld *SMRecord, infrastructure string) error { //do zmiany

	// sm_json, err := json.Marshal(smRecord)
	// if err != nil {
	// 	return err
	// }
	// logger.Debug(string(sm_json))
	// data := url.Values{"parameters": {string(sm_json)}, "infrastructure": {infrastructure}}

	//----
	data := url.Values{"parameters": {sm_record_marshal(smRecord, smRecordOld)}, "infrastructure": {infrastructure}}
	//----

	urlString := fmt.Sprintf("%v://%v/simulation_managers/%v", emc.scheme, emc.experimentManagerAddress, smRecord.ID)

	request, err := http.NewRequest("PUT", urlString, strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	request.SetBasicAuth(emc.login, emc.password)

	resp, err := emc.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	} else {
		logger.Info("%v Status code: %v", smRecord.GetIDs(), resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		logger.Debug("%s", body)
		return errors.New("Update failed")
	}
	return nil
}

func escape(s string) string {
	s = strings.Replace(s, "\n", "\\n", -1)
	s = strings.Replace(s, "\r", "\\r", -1)
	s = strings.Replace(s, "\t", "\\t", -1)
	s = strings.Replace(s, `'`, `\'`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)

	return s
}
