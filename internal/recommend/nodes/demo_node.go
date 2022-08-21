package nodes

import (
	"fmt"
	"grahp-impl/internal/configs"
	"grahp-impl/internal/recommend"
	"time"
)

type DemoNode struct {
}

func (m *DemoNode) Init(config *configs.Config) bool {
	return true
}

func (m *DemoNode) Handle(ctx *recommend.Ctx, xhg *recommend.Exchange, attr recommend.Attribute) {
	if attr.Count < 0 {
		panic("count < 0")
	}
	time.Sleep(time.Duration(attr.Count) * time.Second)
	fmt.Printf("[%s] sleep %d second\n", attr.Name, attr.Count)
}
