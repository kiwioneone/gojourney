package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type ThreadListSync struct {
	Data *discord.ThreadListSyncEventFields `json:"d"`
}

func NewThreadListSync(rest *rest.Client, data []byte) (*ThreadListSync, error) {
	pk := new(ThreadListSync)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
