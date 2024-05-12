package cmd

import (
	"encoding/xml"
	"fmt"
	"mon-dell-me4012/config"

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
	Name string `json:"{#NAME}"`
}

type totalSensorStatuses struct {
	Data []any `json:"data"`
}

func DiscoverySensorStatus(session *ssh.Session, c *config.Config) (totalSensorStatuses, error) {
	v, err := getRawData(session, "SensorStatus", "show sensor-status", c.Redis.SSHBlockExpireKeyTime, c.Redis.DataExpireKeyTime)
	if err != nil {
		return totalSensorStatuses{}, fmt.Errorf("%s", err)
	}

	var XMLData SensorStatus = SensorStatus{}
	err = xml.Unmarshal([]byte(v), &XMLData)
	if err != nil {
		return totalSensorStatuses{}, fmt.Errorf("error XML unmarshal SensorStatus (discovery): %v", err)
	}

	var fakeSingleEntity fakeSingleEntity
	var totalSensorStatuses totalSensorStatuses
	var singleSensorStatus singleSensorStatus
	for _, sensorStatuse := range XMLData.OBJECT {
		sensorStatuse.Name = sensorStatuse.PROPERTY[0].Text
		if sensorStatuse.Name == "Success" {
			continue
		}

		singleSensorStatus.Name = sensorStatuse.PROPERTY[0].Text

		totalSensorStatuses.Data = append(totalSensorStatuses.Data, singleSensorStatus)
	}

	totalSensorStatuses.Data = append(totalSensorStatuses.Data, fakeSingleEntity)

	return totalSensorStatuses, nil
}

func GetValuesBySensorStatus(session *ssh.Session, c *config.Config, sensorStatusName, param string) (string, error) {
	result := map[string]any{}

	v, err := getRawData(session, "SensorStatus", "show sensor-status", c.Redis.SSHBlockExpireKeyTime, c.Redis.DataExpireKeyTime)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	var XMLData SensorStatus = SensorStatus{}
	err = xml.Unmarshal([]byte(v), &XMLData)
	if err != nil {
		return "", fmt.Errorf("error XML unmarshal SensorStatus (specific): %v", err)
	}

	for _, i := range XMLData.OBJECT {
		if i.PROPERTY[0].Text == sensorStatusName {
			result["Health"] = i.PROPERTY[8].Text
		}
	}

	return fmt.Sprintf("%v", result[param]), nil
}
