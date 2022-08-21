package recommend

import (
	"fmt"
	"grahp-impl/pkg/schema"
	"runtime"
	"sync"
)

type innerPanic string

var WorkflowExecutor workflowExecutor

type workflowExecutor struct {
}

func (m *workflowExecutor) Run(workflow *Workflow) {

	wg := sync.WaitGroup{}
	wg.Add(len(workflow.Sections))
	for _, section := range workflow.Sections {
		go func(section *Section, workflow *Workflow) {

			defer func() {

				recoverFromPanic := false
				if r := recover(); r != nil {
					recoverFromPanic = true
					fmt.Println("Panic:", r)
					if _, ok := r.(innerPanic); !ok {
						buf := make([]byte, 10240)
						l := runtime.Stack(buf, false)
						fmt.Println(string(buf[:l]))
					}
				}

				workflowNodeStatus := WorkflowStatus_OK
				if recoverFromPanic {
					workflowNodeStatus = WorkflowStatus_Panic
				}
				section.status = workflowNodeStatus

				wg.Done()
			}()

			m.executeSection(section, workflow)

		}(section, workflow)
	}
	wg.Wait()

}

func (m *workflowExecutor) executeSection(section *Section, workflow *Workflow) {

	if section.Realtime != nil && section.Realtime.Enable == RealtimeStatus_ON {
		err := InitRealtime(section)
		if err != nil {
			section.Realtime.Enable = RealtimeStatus_OFF
		}
	}

	// Parallelly execute all nodes, and wait until all nodes are finished
	wg := sync.WaitGroup{}
	wg.Add(len(section.Nodes))
	for _, node := range section.Nodes {
		go func(section *Section, node *Node) {

			defer func() {

				recoverFromPanic := false
				if r := recover(); r != nil {
					recoverFromPanic = true
					fmt.Println("Panic:", r)
					if _, ok := r.(innerPanic); !ok {
						buf := make([]byte, 10240)
						l := runtime.Stack(buf, false)
						fmt.Println(string(buf[:l]))
					}
				}

				workflowNodeStatus := WorkflowStatus_OK
				if recoverFromPanic {
					workflowNodeStatus = WorkflowStatus_Panic
				}

				node.notify(workflowNodeStatus)

				wg.Done()
			}()

			executableNode := ExecutableNodeManager.Get(node.ClassName)
			if executableNode == nil {
				panic(innerPanic(fmt.Sprintf("invalid executable node class[%s] in workflow [%s]", node.ClassName, workflow.Name)))
			}

			panicHappened := false
			panicName := ""
			nodeInput := make([]*schema.Unit, 0, 500)
			for _, inputNodeName := range node.InputNodes {

				inputNode, ok := section.Nodes[inputNodeName]
				if !ok {
					panic(innerPanic(fmt.Sprintf("invalid node name[%s] in workflow [%s]", inputNodeName, workflow.Name)))
				}

				inputNode.wait()

				// record panic info but to wait other nodes to be finished
				if inputNode.Status == WorkflowStatus_Panic {
					panicHappened = true
					panicName = inputNode.Name
				}

				nodeInput = append(nodeInput, inputNode.Exchange.output...)
			}
			if panicHappened {
				panic(innerPanic(fmt.Sprintf("input node [%s] of node [%s] is panic in workflow [%s]", panicName, node.Name, workflow.Name)))
			}
			node.Exchange.input = nodeInput

			//
			executableNode.Handle((*Ctx)(section), &node.Exchange, node.Attribute)

		}(section, node)
	}
	wg.Wait()

	endNode, ok := section.Nodes[section.NameOfEndNode]
	if !ok {
		panic(innerPanic(fmt.Sprintf("invalid end node name[%s] in workflow [%s]", endNode.Name, workflow.Name)))
	}
	if endNode.Status == WorkflowStatus_Panic {
		panic(innerPanic(fmt.Sprintf("end node [%s] of section [%s] is panic in workflow [%s]", endNode.Name, section.Name, workflow.Name)))
	}

	section.status = endNode.Status
	section.result = endNode.Exchange.output

	if section.Realtime != nil && section.Realtime.Enable == RealtimeStatus_ON {
		FinishRealtime(section)
	} else {
		if len(section.result) > section.Capacity {
			section.result = section.result[:section.Capacity]
		}
		section.Capacity = len(section.result)
	}

}
