package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type AutoModerationRuleCreate struct {
	Data *discord.AutoModerationRule `json:"d"`
}

func NewAutoModerationRuleCreate(rest *rest.Client, data []byte) (*AutoModerationRuleCreate, error) {
	pk := new(AutoModerationRuleCreate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
