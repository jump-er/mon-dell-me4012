package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func ExecCommandOnDevice(session *ssh.Session, command string) (bytes.Buffer, error) {
	var buff bytes.Buffer
	session.Stdout = &buff
	if err := session.Run(command); err != nil {
		return buff, fmt.Errorf("error command exec: %v", err)
	}

	return buff, nil
}

func GetAndPrintMetricData(metricData interface{}) {
	jsonResult, err := json.Marshal(metricData)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(jsonResult))
}
