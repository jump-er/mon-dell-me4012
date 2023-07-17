package cmd

import (
	"encoding/xml"
	"fmt"

	"golang.org/x/crypto/ssh"
)

type SensorStatus struct {
	XMLName xml.Name `xml:"RESPONSE"`
	Text    string   `xml:",chardata"`
	VERSION string   `xml:"VERSION,attr"`
	REQUEST string   `xml:"REQUEST,attr"`
	COMP    []struct {
		Text string `xml:",chardata"`
		G    string `xml:"G,attr"`
		P    string `xml:"P,attr"`
	} `xml:"COMP"`
	OBJECT []struct {
		Text     string `xml:",chardata"`
		Basetype string `xml:"basetype,attr"`
		Name     string `xml:"name,attr"`
		Oid      string `xml:"oid,attr"`
		Format   string `xml:"format,attr"`
		PROPERTY []struct {
			Text        string `xml:",chardata"`
			Name        string `xml:"name,attr"`
			Type        string `xml:"type,attr"`
			Size        string `xml:"size,attr"`
			Draw        string `xml:"draw,attr"`
			Sort        string `xml:"sort,attr"`
			DisplayName string `xml:"display-name,attr"`
			Key         string `xml:"key,attr"`
		} `xml:"PROPERTY"`
	} `xml:"OBJECT"`
}

type singleSensorStatus struct {
	Name   string `json:"{#NAME}"`
	Health string `json:"{#HEALTH}"`
}

type totalSensorStatuses struct {
	Data []singleSensorStatus `json:"data"`
}

func GetSensorStatus(session *ssh.Session) (totalSensorStatuses, error) {
	buff, err := ExecCommandOnDevice(session, "show sensor-status")
	if err != nil {
		return totalSensorStatuses{}, fmt.Errorf("%v", err)
	}

	var sensorStatuse SensorStatus
	err = xml.Unmarshal([]byte(buff.String()), &sensorStatuse)
	if err != nil {
		return totalSensorStatuses{}, fmt.Errorf("error XML unmarshal: %v", err)
	}

	var totalSensorStatuses totalSensorStatuses
	var singleSensorStatus singleSensorStatus
	for _, sensorStatuse := range sensorStatuse.OBJECT {
		sensorStatuse.Name = sensorStatuse.PROPERTY[0].Text
		if sensorStatuse.Name == "Success" {
			continue
		}

		singleSensorStatus.Name = sensorStatuse.PROPERTY[0].Text

		singleSensorStatus.Health = sensorStatuse.PROPERTY[8].Text

		totalSensorStatuses.Data = append(totalSensorStatuses.Data, singleSensorStatus)
	}

	return totalSensorStatuses, nil
}
