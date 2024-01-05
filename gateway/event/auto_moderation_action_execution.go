package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type AutoModerationActionExecution struct {
	Data *discord.AutoModerationActionExecutionEventFields `json:"d"`
}

func NewAutoModerationActionExecution(rest *rest.Client, data []byte) (*AutoModerationActionExecution, error) {
	pk := new(AutoModerationActionExecution)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
