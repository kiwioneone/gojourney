package gateway

import (
	"github.com/kiwioneone/gojourney/gateway/event"
)

type InteractionCreateHandler struct{}

func (_ *InteractionCreateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewInteractionCreate(s.rest, data)

	if err != nil {
		return
	}

	s.Publish(event.EventInteractionCreate, ev.Data)
}
