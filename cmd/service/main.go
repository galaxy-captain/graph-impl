package main

import (
	"fmt"
	"grahp-impl/pkg/yaml"
	"os"
)

type Config struct {
	Service ServiceConfig `yaml:"service"`
}

type ServiceConfig struct {
	Port    string `yaml:"port"`
	Timeout int    `yaml:"timeout"`
}

func main() {

	env := os.Getenv("ENV")
	serviceScenario := os.Getenv("SRV_SCENARIO")
	serviceType := os.Getenv("SRV_TYPE")
	port := os.Getenv("PORT")
	fmt.Printf("\n"+
		"ENV: %s\n"+
		"SRV_SCENARIO: %s\n"+
		"SRV_TYPE: %s\n"+
		"PORT: %s\n"+
		"\n", env, serviceScenario, serviceType, port)

	var err error

	filepath := fmt.Sprintf("./configs/%s/service-%s-%s.yaml", serviceScenario, env, serviceType)
	config := &Config{}
	err = yaml.LoadYamlFile(filepath, config)
	if err != nil {
		panic(err)
	}

	fmt.Println(config)
}
