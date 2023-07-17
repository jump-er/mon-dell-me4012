package main

import (
	"flag"
	"fmt"
	"mon-dell-me4012/cmd"
	"mon-dell-me4012/config"
	"strconv"

	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh"
)

func main() {
	config.LoggerInit()

	c, err := config.ConfigInit()
	if err != nil {
		log.Fatalf("fail get config: %v", err)
	}

	hostKeyCallback, err := config.SetInsecureIgnoreHostKeyOption(c)
	if err != nil {
		log.Fatal(err)
	}

	s := &ssh.ClientConfig{
		User: c.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.SSH.Password),
		},
		HostKeyCallback: hostKeyCallback,
	}

	conn, err := ssh.Dial("tcp", c.SSH.Host+":"+strconv.Itoa(c.SSH.Port), s)
	if err != nil {
		log.Fatalf("fail SSH dial: %v", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("fail session: %v", err)
	}
	defer session.Close()

	useMetric := flag.String("metric_name", "", "metric name for getting data")
	flag.Parse()

	switch *useMetric {
	case "VolumeStatistics":
		metricData, err := cmd.GetVolumeStatistics(session)
		if err != nil {
			log.Error(err)
		}
		cmd.GetAndPrintMetricData(metricData)
	case "DiskGroups":
		metricData, err := cmd.GetDiskGroups(session)
		if err != nil {
			log.Error(err)
		}
		cmd.GetAndPrintMetricData(metricData)
	case "PowerSupplies":
		metricData, err := cmd.GetPowerSupplies(session)
		if err != nil {
			log.Error(err)
		}
		cmd.GetAndPrintMetricData(metricData)
	case "SensorStatus":
		metricData, err := cmd.GetSensorStatus(session)
		if err != nil {
			log.Error(err)
		}
		cmd.GetAndPrintMetricData(metricData)
	default:
		fmt.Println("No valid metric name")
	}
}
