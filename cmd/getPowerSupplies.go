package cmd

import (
	"encoding/xml"
	"fmt"
	"mon-dell-me4012/config"
	"mon-dell-me4012/rds"
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
	Name string `json:"{#NAME}"`
}

type totalPowerSupplies struct {
	Data []any `json:"data"`
}

func DiscoveryPowerSupplies(session *ssh.Session, c *config.Config) (totalPowerSupplies, error) {
	v, err := getRawData(session, rds.R, "PowerSupplies", "show power-supplies", c.Redis.SSHBlockExpireKeyTime, c.Redis.DataExpireKeyTime)
	if err != nil {
		return totalPowerSupplies{}, fmt.Errorf("%s", err)
	}

	var XMLData PowerSupplies = PowerSupplies{}
	err = xml.Unmarshal([]byte(v), &XMLData)
	if err != nil {
		return totalPowerSupplies{}, fmt.Errorf("error XML unmarshal PowerSupplies (discovery): %v", err)
	}

	var fakeSingleEntity fakeSingleEntity
	var totalPowerSupplies totalPowerSupplies
	var singlePowerSupplie singlePowerSupplie
	for _, powerSupplie := range XMLData.OBJECT {
		powerSupplie.Name = powerSupplie.PROPERTY[0].Text
		if powerSupplie.Name == "Success" || strings.Contains(powerSupplie.Name, "fan") {
			continue
		}

		singlePowerSupplie.Name = powerSupplie.PROPERTY[0].Text

		totalPowerSupplies.Data = append(totalPowerSupplies.Data, singlePowerSupplie)
	}

	totalPowerSupplies.Data = append(totalPowerSupplies.Data, fakeSingleEntity)

	return totalPowerSupplies, nil
}

func GetValuesByPowerSupplie(session *ssh.Session, c *config.Config, powerSupplieName, param string) (string, error) {
	result := map[string]any{}

	v, err := getRawData(session, rds.R, "PowerSupplies", "show power-supplies", c.Redis.SSHBlockExpireKeyTime, c.Redis.DataExpireKeyTime)
	if err != nil {
		return "", fmt.Errorf("%s", err)
	}

	var XMLData PowerSupplies = PowerSupplies{}
	err = xml.Unmarshal([]byte(v), &XMLData)
	if err != nil {
		return "", fmt.Errorf("error XML unmarshal PowerSupplies (specific): %v", err)
	}

	for _, i := range XMLData.OBJECT {
		if i.PROPERTY[0].Text == powerSupplieName {
			result["Name"] = i.PROPERTY[0].Text

			result["Health"] = i.PROPERTY[27].Text

		}
	}

	return fmt.Sprintf("%v", result[param]), nil
}
