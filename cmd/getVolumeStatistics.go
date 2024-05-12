package cmd

import (
	"encoding/xml"
	"fmt"
	"mon-dell-me4012/config"
	"mon-dell-me4012/rds"

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
}

type totalVolumes struct {
	Data []any `json:"data"`
}

func DiscoveryVolumeStatistics(session *ssh.Session, c *config.Config) (totalVolumes, error) {
	v, err := getRawData(session, rds.R, "VolumeStatistics", "show volume-statistics", c.Redis.SSHBlockExpireKeyTime, c.Redis.DataExpireKeyTime)
	if err != nil {
		return totalVolumes{}, fmt.Errorf("%s", err)
	}

	var XMLData VolumeStatistic = VolumeStatistic{}
	err = xml.Unmarshal([]byte(v), &XMLData)
	if err != nil {
		return totalVolumes{}, fmt.Errorf("error XML unmarshal VolumeStatistics (discovery): %v", err)
	}

	var fakeSingleEntity fakeSingleEntity
	var totalVolumes totalVolumes
	var singleVolume singleVolume
	for _, volume := range XMLData.OBJECT {
		singleVolume.Name = volume.PROPERTY[0].Text
		if singleVolume.Name == "Success" {
			continue
		}

		totalVolumes.Data = append(totalVolumes.Data, singleVolume)
	}

	totalVolumes.Data = append(totalVolumes.Data, fakeSingleEntity)

	return totalVolumes, nil
}

func GetValuesByVolume(session *ssh.Session, c *config.Config, volumeName, param string) (string, error) {
	result := map[string]any{}

	v, err := getRawData(session, rds.R, "VolumeStatistics", "show volume-statistics", c.Redis.SSHBlockExpireKeyTime, c.Redis.DataExpireKeyTime)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	var XMLData VolumeStatistic = VolumeStatistic{}
	err = xml.Unmarshal([]byte(v), &XMLData)
	if err != nil {
		return "", fmt.Errorf("error XML unmarshal VolumeStatistics (specific): %v", err)
	}

	for _, i := range XMLData.OBJECT {
		if i.PROPERTY[0].Text == volumeName {
			result["BPS"] = i.PROPERTY[3].Text

			result["IOPS"] = i.PROPERTY[4].Text
		}
	}

	return fmt.Sprintf("%v", result[param]), nil
}
