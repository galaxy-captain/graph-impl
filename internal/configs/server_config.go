package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port    string `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

func LoadYamlFile(filepath string) (*Config, error) {

	// read data from file
	serviceConfigYaml, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	serviceConfigYamlStr := string(serviceConfigYaml)

	// apply env variables to data
	allEnvVars := os.Environ()
	for _, env := range allEnvVars {
		kv := strings.Split(env, "=")
		if len(kv) != 2 {
			continue
		}
		key := fmt.Sprintf("env:%s", kv[0])
		value := kv[1]
		serviceConfigYamlStr = strings.ReplaceAll(serviceConfigYamlStr, key, value)
	}

	// deserialize data to object
	config := new(Config)
	err = yaml.Unmarshal([]byte(serviceConfigYamlStr), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
