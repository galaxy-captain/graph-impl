package nodes

import (
	"grahp-impl/internal/configs"
	"grahp-impl/internal/recommend"
	"net/http"
)

type RecallNode struct {
}

func (m *RecallNode) Init(config *configs.Config) bool {
	return true
}

func (m *RecallNode) Handle(ctx *recommend.Ctx, xhg *recommend.Exchange, attr recommend.Attribute) {

	_, err := http.NewRequest(http.MethodPost, "", nil)
	if err != nil {
		return
	}

}
