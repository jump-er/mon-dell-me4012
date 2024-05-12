package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"gopkg.in/yaml.v2"
)

var configFilePath string = "./config.yaml"

type Config struct {
	InsecureIgnoreHostKey bool   `yaml:"insecure_ignore_host_key"`
	KnownhostsFilePath    string `yaml:"knownhosts_file_path"`
	SSH                   struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
	} `yaml:"ssh"`
	Redis struct {
		Host                  string `yaml:"host"`
		Port                  int    `yaml:"port"`
		SSHBlockExpireKeyTime int    `yaml:"ssh_block_expire_key_time"`
		DataExpireKeyTime     int    `yaml:"data_expire_key_time"`
	} `yaml:"redis"`
}

func ConfigInit() (*Config, error) {
	f, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("fail reading config file")
	}

	var c Config
	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return nil, fmt.Errorf("fail unmarshaling config")
	}

	return &c, nil
}

func LoggerInit() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	file, err := os.OpenFile("./mon-dell-me4012.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Warn("Failed to log to file, using default stderr")
	}
}

func SetInsecureIgnoreHostKeyOption(config *Config) (ssh.HostKeyCallback, error) {
	if config.InsecureIgnoreHostKey {
		return ssh.InsecureIgnoreHostKey(), nil
	}

	hostKeyCallback, err := knownhosts.New(config.KnownhostsFilePath)
	if err != nil {
		return nil, fmt.Errorf("fail get knownhosts file: %v", err)
	}

	return hostKeyCallback, nil
}
