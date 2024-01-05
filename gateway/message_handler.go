package gateway

import (
	"github.com/kiwioneone/gojourney/gateway/event"
)

type MessageCreateHandler struct{}

func (_ *MessageCreateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewMessageCreate(s.rest, data)

	if err != nil {
		return
	}

	s.Publish(event.EventMessageCreate, ev.Data)
}

type MessageUpdateHandler struct{}

func (_ *MessageUpdateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewMessageCreate(s.rest, data)

	if err != nil {
		return
	}

	s.Publish(event.EventMessageUpdate, ev.Data)
}

type MessageDeleteHandler struct{}

func (_ *MessageDeleteHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewMessageDelete(s.rest, data)

	if err != nil {
		return
	}

	s.Publish(event.EventMessageDelete, ev.Data)
}
