package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mon-dell-me4012/rds"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type fakeSingleEntity struct {
	Name string `json:"{#}"`
}

func execCommandOnDevice(session *ssh.Session, command string) (bytes.Buffer, error) {
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

func getRawData(session *ssh.Session, rc *rds.RedisCache, entity, SSHcommand string, EXs, EXd int) (string, error) {
	s, _ := rc.Get(rc.Client, entity)
	if s == "" {
		buff, err := execCommandOnDevice(session, SSHcommand)
		if err != nil {
			return "", fmt.Errorf("error exec '%s' SSH commend: %s", SSHcommand, err)
		}

		err = rc.Set(rc.Client, entity, buff.String(), EXd)
		if err != nil {
			return "", fmt.Errorf("error setting %s data to Redis: %s", entity, err)
		}

		err = rc.Set(rc.Client, fmt.Sprintf("SSHConnectionBlockFrom%s", entity), "yes", EXs)
		if err != nil {
			return "", fmt.Errorf("error setting SSHConnectionBlock data to Redis: %s", err)
		}

		s, err = rc.Get(rc.Client, entity)
		if err != nil {
			return "", fmt.Errorf("error getting %s data from Redis after SET: %s", entity, err)
		}
	}

	return s, nil
}
