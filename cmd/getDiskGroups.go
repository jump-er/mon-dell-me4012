package cmd

import (
	"encoding/xml"
	"fmt"
	"mon-dell-me4012/config"
	"mon-dell-me4012/rds"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type DiskGroups struct {
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
			Units       string `xml:"units,attr"`
			Key         string `xml:"key,attr"`
			Blocksize   string `xml:"blocksize,attr"`
		} `xml:"PROPERTY"`
	} `xml:"OBJECT"`
}

type singleDiskGroup struct {
	Name string `json:"{#NAME}"`
}

type totalDiskGroups struct {
	Data []any `json:"data"`
}

func DiscoveryDiskGroups(session *ssh.Session, c *config.Config) (totalDiskGroups, error) {
	v, err := getRawData(session, rds.R, "DiskGroups", "show disk-groups", c.Redis.SSHBlockExpireKeyTime, c.Redis.DataExpireKeyTime)
	if err != nil {
		return totalDiskGroups{}, fmt.Errorf("%s", err)
	}

	var XMLData DiskGroups = DiskGroups{}
	err = xml.Unmarshal([]byte(v), &XMLData)
	if err != nil {
		return totalDiskGroups{}, fmt.Errorf("error XML unmarshal DiskGroups (discovery): %v", err)
	}

	var fakeSingleEntity fakeSingleEntity
	var totalDiskGroups totalDiskGroups
	var singleDiskGroup singleDiskGroup
	for _, diskGroup := range XMLData.OBJECT {
		diskGroup.Name = diskGroup.PROPERTY[0].Text
		if diskGroup.Name == "Success" {
			continue
		}

		singleDiskGroup.Name = diskGroup.PROPERTY[0].Text

		totalDiskGroups.Data = append(totalDiskGroups.Data, singleDiskGroup)
	}

	totalDiskGroups.Data = append(totalDiskGroups.Data, fakeSingleEntity)

	return totalDiskGroups, nil
}

func GetValuesByDiskGroup(session *ssh.Session, c *config.Config, diskGroupName, param string) (string, error) {
	result := map[string]any{}

	v, err := getRawData(session, rds.R, "DiskGroups", "show disk-groups", c.Redis.SSHBlockExpireKeyTime, c.Redis.DataExpireKeyTime)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	var XMLData DiskGroups = DiskGroups{}
	err = xml.Unmarshal([]byte(v), &XMLData)
	if err != nil {
		return "", fmt.Errorf("error XML unmarshal DiskGroups (specific): %v", err)
	}

	for _, i := range XMLData.OBJECT {
		if i.PROPERTY[0].Text == diskGroupName {
			s, err := strconv.Atoi(i.PROPERTY[3].Text)
			if err != nil {
				return "", fmt.Errorf("%s", err)
			}
			result["Size"] = (s / 2) * 1000

			s, err = strconv.Atoi(i.PROPERTY[5].Text)
			if err != nil {
				return "", fmt.Errorf("%s", err)
			}
			result["Free"] = (s / 2) * 1000

			result["Status"] = i.PROPERTY[28].Text

			result["Health"] = i.PROPERTY[77].Text
		}
	}

	return fmt.Sprintf("%v", result[param]), nil
}
