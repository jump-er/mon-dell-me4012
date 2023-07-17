package cmd

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type VolumeStatistic struct {
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
			BPS         string `xml:"size,attr"`
			Draw        string `xml:"draw,attr"`
			Sort        string `xml:"sort,attr"`
			DisplayName string `xml:"display-name,attr"`
			Key         string `xml:"key,attr"`
			Units       string `xml:"units,attr"`
		} `xml:"PROPERTY"`
	} `xml:"OBJECT"`
}

type singleVolume struct {
	Name string `json:"{#NAME}"`
	BPS  int    `json:"{#BPS}"`
	IOPS int    `json:"{#IOPS}"`
}

type totalVolumes struct {
	Data []singleVolume `json:"data"`
}

func GetVolumeStatistics(session *ssh.Session) (totalVolumes, error) {
	buff, err := ExecCommandOnDevice(session, "show volume-statistics")
	if err != nil {
		return totalVolumes{}, fmt.Errorf("%v", err)
	}

	var volumeStatistic VolumeStatistic
	err = xml.Unmarshal([]byte(buff.String()), &volumeStatistic)
	if err != nil {
		return totalVolumes{}, fmt.Errorf("error XML unmarshal: %v", err)
	}

	var totalVolumes totalVolumes
	var singleVolume singleVolume
	for _, volume := range volumeStatistic.OBJECT {
		singleVolume.Name = volume.PROPERTY[0].Text
		if singleVolume.Name == "Success" {
			continue
		}

		singleVolume.BPS, err = strconv.Atoi(volume.PROPERTY[3].Text)
		if err != nil {
			return totalVolumes, fmt.Errorf("error convert BPS count to int: %v", err)
		}

		singleVolume.IOPS, err = strconv.Atoi(volume.PROPERTY[4].Text)
		if err != nil {
			return totalVolumes, fmt.Errorf("error convert IOPS count to int: %v", err)
		}

		totalVolumes.Data = append(totalVolumes.Data, singleVolume)
	}

	return totalVolumes, nil
}
