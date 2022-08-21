package recommend

import (
	"fmt"
	"grahp-impl/internal/configs"
)

type Ctx Section

type ExecutableNode interface {
	Init(config *configs.Config) bool
	Handle(ctx *Ctx, exchange *Exchange, attr Attribute)
}

var ExecutableNodeManager executableNodeManager

type executableNodeManager struct {
	nodeMap map[string]ExecutableNode
}

func (m *executableNodeManager) Add(name string, executableNode ExecutableNode) {
	if m.nodeMap == nil {
		m.nodeMap = make(map[string]ExecutableNode)
	}
	m.nodeMap[name] = executableNode
}

func (m *executableNodeManager) Get(name string) ExecutableNode {
	if executableNode, ok := m.nodeMap[name]; !ok {
		return nil
	} else {
		return executableNode
	}
}

func (m *executableNodeManager) Init(config *configs.Config) error {
	for name, node := range m.nodeMap {
		if node.Init(config) == false {
			return fmt.Errorf("failed to init node [%s]", name)
		}
	}
	return nil
}
