package nodes

import (
	"grahp-impl/internal/configs"
	"grahp-impl/internal/recommend"
)

func RegisterAllNodes(config *configs.Config) error {

	recommend.ExecutableNodeManager.Add("DemoNode", new(DemoNode))

	err := recommend.ExecutableNodeManager.Init(config)
	if err != nil {
		return err
	}

	return nil
}
