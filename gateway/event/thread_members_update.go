package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type ThreadMembersUpdate struct {
	Data *discord.ThreadMembersUpdateEventFields `json:"d"`
}

func NewThreadMembersUpdate(rest *rest.Client, data []byte) (*ThreadMembersUpdate, error) {
	pk := new(ThreadMembersUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
