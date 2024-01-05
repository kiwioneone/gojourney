package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type ThreadMemberUpdate struct {
	Data *discord.ThreadMember `json:"d"`
}

func NewThreadMemberUpdate(rest *rest.Client, data []byte) (*ThreadMemberUpdate, error) {
	pk := new(ThreadMemberUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
