package cmd

import (
	"encoding/xml"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

type PowerSupplies struct {
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
			Deprecated  string `xml:"deprecated,attr"`
		} `xml:"PROPERTY"`
	} `xml:"OBJECT"`
}

type singlePowerSupplie struct {
	Name   string `json:"{#NAME}"`
	Health string `json:"{#HEALTH}"`
}

type totalPowerSupplies struct {
	Data []singlePowerSupplie `json:"data"`
}

func GetPowerSupplies(session *ssh.Session) (totalPowerSupplies, error) {
	buff, err := ExecCommandOnDevice(session, "show power-supplies")
	if err != nil {
		return totalPowerSupplies{}, fmt.Errorf("%v", err)
	}

	var powerSupplies PowerSupplies
	err = xml.Unmarshal([]byte(buff.String()), &powerSupplies)
	if err != nil {
		return totalPowerSupplies{}, fmt.Errorf("error XML unmarshal: %v", err)
	}

	var totalPowerSupplies totalPowerSupplies
	var singlePowerSupplie singlePowerSupplie
	for _, powerSupplie := range powerSupplies.OBJECT {
		powerSupplie.Name = powerSupplie.PROPERTY[0].Text
		if powerSupplie.Name == "Success" || strings.Contains(powerSupplie.Name, "fan") {
			continue
		}

		singlePowerSupplie.Name = powerSupplie.PROPERTY[0].Text

		singlePowerSupplie.Health = powerSupplie.PROPERTY[27].Text

		totalPowerSupplies.Data = append(totalPowerSupplies.Data, singlePowerSupplie)
	}

	return totalPowerSupplies, nil
}
