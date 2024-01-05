package gateway

import (
	"github.com/kiwioneone/gojourney/gateway/event"
)

type ApplicationCommandPermissionsUpdateHandler struct{}

func (_ *ApplicationCommandPermissionsUpdateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewApplicationCommandPermissionsUpdate(s.rest, data)

	if err != nil {
		return
	}

	s.Publish(event.EventApplicationCommandPermissionsUpdate, ev.Data)
}
