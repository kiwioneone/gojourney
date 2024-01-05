package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type AutoModerationRuleDelete struct {
	Data *discord.AutoModerationRule `json:"d"`
}

func NewAutoModerationRuleDelete(rest *rest.Client, data []byte) (*AutoModerationRuleDelete, error) {
	pk := new(AutoModerationRuleDelete)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
