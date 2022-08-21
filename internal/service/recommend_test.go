package service

import (
	"fmt"
	"grahp-impl/internal/configs"
	"grahp-impl/internal/recommend"
	"grahp-impl/internal/recommend/nodes"
)

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func init() {

	config := new(configs.Config)

	err := nodes.RegisterAllNodes(config)
	if err != nil {
		fmt.Println(err)
	}
}

func TestBuildWorkflow(t *testing.T) {

	rawData, err := ioutil.ReadFile("./recommend_testdata.json")
	if err != nil {
		t.Error(err)
	}

	workflowMap := recommend.WorkflowMap{}
	err = json.Unmarshal(rawData, &workflowMap)
	if err != nil {
		t.Error(err)
	}

	for workflowName, workflowConfig := range workflowMap.Workflows {

		workflow, err := recommend.BuildWorkflow(workflowName, &workflowConfig)
		if err != nil {
			t.Error(err)
		}

		recommend.WorkflowExecutor.Run(workflow)
	}

}
