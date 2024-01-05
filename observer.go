package gojourney

import (
	"github.com/kiwioneone/gojourney/discord"
)

type InteractionObserver interface {
	Observe(interaction *discord.Interaction)
}

type MessageObserver interface {
	Observe(msg *discord.Message)
}
