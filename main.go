package main

import (
	"flag"
	"fmt"
	"mon-dell-me4012/cmd"
	"mon-dell-me4012/config"
	"mon-dell-me4012/rds"
	"strconv"

	log "github.com/sirupsen/logrus"

	"golang.org/x/crypto/ssh"
)

// SSHConnectionBlockFromVolumeStatistics не устанавливается в "yes"

func main() {
	config.LoggerInit()

	c, err := config.ConfigInit()
	if err != nil {
		log.Fatal(err)
	}

	rds.R = rds.RedisInit(c)

	discovery := flag.Bool("discovery", false, "discovery behavior")
	discoveryName := flag.String("discovery_name", "", "name for specific discovery")
	metricGroup := flag.String("metric_group", "", "metric group")
	entityName := flag.String("entity_name", "", "entity name")
	metricName := flag.String("metric_name", "", "metric name")
	flag.Parse()

	var session *ssh.Session = &ssh.Session{}
	if !isSSHConnectionBlocked(getEntityName(discoveryName, metricGroup)) {
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

		session, err = conn.NewSession()
		if err != nil {
			log.Fatalf("fail session: %v", err)
		}
		defer session.Close()
	}

	if *discovery {
		switch *discoveryName {
		case "VolumeStatistics":
			metricData, err := cmd.DiscoveryVolumeStatistics(session, c)
			if err != nil {
				log.Error(err)
			}
			cmd.GetAndPrintMetricData(metricData)
			return
		case "DiskGroups":
			metricData, err := cmd.DiscoveryDiskGroups(session, c)
			if err != nil {
				log.Error(err)
			}
			cmd.GetAndPrintMetricData(metricData)
			return
		case "PowerSupplies":
			metricData, err := cmd.DiscoveryPowerSupplies(session, c)
			if err != nil {
				log.Error(err)
			}
			cmd.GetAndPrintMetricData(metricData)
			return
		case "SensorStatus":
			metricData, err := cmd.DiscoverySensorStatus(session, c)
			if err != nil {
				log.Error(err)
			}
			cmd.GetAndPrintMetricData(metricData)
			return
		default:
			fmt.Println("No valid discovery name")
			return
		}
	}

	switch *metricGroup {
	case "DiskGroups":
		v, err := cmd.GetValuesByDiskGroup(session, c, *entityName, *metricName)
		if err != nil {
			log.Error(err)
		}
		fmt.Print(v)
	case "PowerSupplies":
		v, err := cmd.GetValuesByPowerSupplie(session, c, *entityName, *metricName)
		if err != nil {
			log.Error(err)
		}
		fmt.Print(v)
	case "SensorStatus":
		v, err := cmd.GetValuesBySensorStatus(session, c, *entityName, *metricName)
		if err != nil {
			log.Error(err)
		}
		fmt.Print(v)
	case "VolumeStatistics":
		v, err := cmd.GetValuesByVolume(session, c, *entityName, *metricName)
		if err != nil {
			log.Error(err)
		}
		fmt.Print(v)
	default:
		fmt.Println("No valid metric group, entity name or metric name")
	}
}

func isSSHConnectionBlocked(entity string) bool {
	SSHConnectionBlock, err := rds.RedisGet(rds.R, fmt.Sprintf("SSHConnectionBlockFrom%s", entity))
	if err != nil {
		log.Errorf("%s", err)
	}

	return SSHConnectionBlock == "yes"
}

func getEntityName(s1, s2 *string) string {
	if *s1 != "" {
		return *s1
	}

	return *s2
}
