package cmd

import (
	"encoding/xml"
	"fmt"
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
	Name   string `json:"{#NAME}"`
	Size   int    `json:"{#SIZE}"`
	Free   int    `json:"{#FREE}"`
	Status string `json:"{#STATUS}"`
	Health string `json:"{#HEALTH}"`
}

type totalDiskGroups struct {
	Data []singleDiskGroup `json:"data"`
}

func GetDiskGroups(session *ssh.Session) (totalDiskGroups, error) {
	buff, err := ExecCommandOnDevice(session, "show disk-groups")
	if err != nil {
		return totalDiskGroups{}, fmt.Errorf("%v", err)
	}

	var diskGroups DiskGroups
	err = xml.Unmarshal([]byte(buff.String()), &diskGroups)
	if err != nil {
		return totalDiskGroups{}, fmt.Errorf("error XML unmarshal: %v", err)
	}

	var totalDiskGroups totalDiskGroups
	var singleDiskGroup singleDiskGroup
	for _, diskGroup := range diskGroups.OBJECT {
		diskGroup.Name = diskGroup.PROPERTY[0].Text
		if diskGroup.Name == "Success" {
			continue
		}

		singleDiskGroup.Name = diskGroup.PROPERTY[0].Text

		i, err := strconv.Atoi(diskGroup.PROPERTY[3].Text)
		if err != nil {
			return totalDiskGroups, fmt.Errorf("error convert Size count to int: %v", err)
		}
		singleDiskGroup.Size = (i / 2) * 1000 //because raid 10

		i, err = strconv.Atoi(diskGroup.PROPERTY[5].Text)
		if err != nil {
			return totalDiskGroups, fmt.Errorf("error convert Free count to int: %v", err)
		}
		singleDiskGroup.Free = (i / 2) * 1000

		singleDiskGroup.Status = diskGroup.PROPERTY[28].Text

		singleDiskGroup.Health = diskGroup.PROPERTY[77].Text

		totalDiskGroups.Data = append(totalDiskGroups.Data, singleDiskGroup)
	}

	return totalDiskGroups, nil
}
