package recommend

import (
	"fmt"
	"grahp-impl/pkg/schema"
	"sync"
)

type WorkflowStatus int32

const WorkflowStatus_None WorkflowStatus = 0
const WorkflowStatus_Panic WorkflowStatus = -1
const WorkflowStatus_OK WorkflowStatus = 1

type Workflow struct {
	Name     string
	Sections map[string]*Section
}

type Section struct {
	Name     string
	Capacity int
	Realtime *Realtime

	Nodes         map[string]*Node
	NameOfEndNode string

	status WorkflowStatus
	result []*schema.Unit
}

func (m *Section) GetResult() []*schema.Unit {
	return m.result
}

func (m *Section) Status() WorkflowStatus {
	return m.status
}

type Node struct {
	Name       string
	ClassName  string
	InputNodes []string

	Exchange Exchange

	Attribute Attribute

	cond   sync.Cond
	Status WorkflowStatus
}

func (m *Node) notify(status WorkflowStatus) {
	m.cond.L.Lock()
	m.Status = status
	m.cond.Broadcast()
	m.cond.L.Unlock()
}

func (m *Node) wait() {
	m.cond.L.Lock()
	for m.Status == WorkflowStatus_None {
		m.cond.Wait()
	}
	m.cond.L.Unlock()
}

func BuildWorkflow(workflowName string, workflowConfig *WorkflowConfig) (*Workflow, error) {

	var err error

	workflow := new(Workflow)
	workflow.Name = workflowName
	workflow.Sections = make(map[string]*Section)
	for sectionName, sectionConfig := range workflowConfig.Sections {

		section := new(Section)
		section.Name = sectionName
		section.Capacity = sectionConfig.Capacity
		section.Nodes = make(map[string]*Node)
		for nodeName, nodeConfig := range sectionConfig.Nodes {

			node := new(Node)
			node.cond.L = new(sync.Mutex)

			node.Name = nodeName
			node.ClassName = nodeConfig.ClassName
			node.InputNodes = nodeConfig.InputNodes

			node.Attribute.Name = nodeName
			node.Attribute.Count = nodeConfig.Count

			section.Nodes[nodeName] = node
		}

		section.NameOfEndNode, err = findEndNodeName(section)
		if err != nil {
			return nil, err
		}

		workflow.Sections[sectionName] = section
	}

	return workflow, nil
}

func findEndNodeName(section *Section) (string, error) {

	nonRootNode := make(map[string]bool)
	for _, node := range section.Nodes {
		for _, inputNodeName := range node.InputNodes {
			nonRootNode[inputNodeName] = true
		}
	}

	rootNodeCount := 0
	rootNode := ""
	for nodeName, _ := range section.Nodes {
		if _, ok := nonRootNode[nodeName]; !ok {
			rootNode = nodeName
			rootNodeCount++
			if rootNodeCount > 1 {
				return "", fmt.Errorf("found more than 1 end node of section[%s]", section.Name)
			}
		}
	}
	if rootNodeCount == 0 {
		return "", fmt.Errorf("not found end node of section[%s]", section.Name)
	}

	return rootNode, nil
}
