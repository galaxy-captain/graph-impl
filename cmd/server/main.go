package main

import (
	"flag"
	"fmt"
	"grahp-impl/internal/configs"
	"os"
)

var ConfigFile = flag.String("config_file", "./configs/%s/service-%s-%s.yaml", "")

func main() {

	flag.Parse()

	var err error

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

	configFile := fmt.Sprintf(*ConfigFile, serviceScenario, env, serviceType)
	config, err := configs.LoadYamlFile(configFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(config)
}
