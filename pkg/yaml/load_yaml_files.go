package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

func LoadYamlFile(filepath string, obj interface{}) error {

	// read data from file
	serviceConfigYaml, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
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
	err = yaml.Unmarshal([]byte(serviceConfigYamlStr), obj)
	if err != nil {
		return err
	}

	return nil
}
