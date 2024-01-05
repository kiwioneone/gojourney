package gateway

import (
	"github.com/kiwioneone/gojourney/gateway/event"
)

type PresenceUpdateHandler struct{}

func (_ *PresenceUpdateHandler) Handle(s *Session, Data []byte) {
	ev, err := event.NewPresenceUpdate(s.rest, Data)

	if err != nil {
		return
	}

	// ToDo : Add a method in state to track presences

	s.Publish(event.EventPresenceUpdate, ev.Data)
}

type GuildCreateHandler struct{}

func (_ *GuildCreateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildCreate(s.rest, data)

	if err != nil {
		return
	}

	s.State().AddGuild(ev.Data)

	s.Publish(event.EventGuildCreate, ev.Data)
}

type GuildUpdateHandler struct{}

func (_ *GuildUpdateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildUpdate(s.rest, data)

	if err != nil {
		return
	}

	s.State().AddGuild(ev.Data)

	s.Publish(event.EventGuildUpdate, ev.Data)
}

type GuildDeleteHandler struct{}

func (_ *GuildDeleteHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildDelete(s.rest, data)

	if err != nil {
		return
	}

	_ = s.State().RemoveGuild(ev.Data)

	s.Publish(event.EventGuildDelete, ev.Data)
}

type GuildBanAddHandler struct{}

func (_ *GuildBanAddHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildBanAdd(s.rest, data)

	if err != nil {
		return
	}

	guild, err := s.State().Guild(ev.Data.GuildId)
	user := ev.Data.User

	if err != nil {
		return
	}

	s.Publish(event.EventGuildBanAdd, guild, user)
}

type GuildBanRemoveHandler struct{}

func (_ *GuildBanRemoveHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildBanRemove(s.rest, data)

	if err != nil {
		return
	}

	guild, err := s.State().Guild(ev.Data.GuildId)
	user := ev.Data.User

	if err != nil {
		return
	}

	s.Publish(event.EventGuildBanRemove, guild, user)
}

type GuildEmojisUpdateHandler struct{}

func (_ *GuildEmojisUpdateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildEmojisUpdate(s.rest, data)

	if err != nil {
		return
	}

	guild, err := s.State().Guild(ev.Data.GuildId)
	if err != nil {
		return
	}

	guild.Emojis = ev.Data.Emojis

	s.Publish(event.EventGuildEmojisUpdate, guild)
}

type GuildStickersUpdateHandler struct{}

func (_ *GuildStickersUpdateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildStickersUpdate(s.rest, data)

	if err != nil {
		return
	}

	// ToDo : Cache stickers?

	s.Publish(event.EventGuildStickersUpdate, ev.Data)
}

type GuildIntegrationsUpdateHandler struct{}

func (_ *GuildIntegrationsUpdateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildIntegrationsUpdate(s.rest, data)

	if err != nil {
		return
	}

	// ToDo : Cache integrations?

	s.Publish(event.EventGuildIntegrationsUpdate, ev.Data)
}

type GuildMemberAddHandler struct{}

func (_ *GuildMemberAddHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildMemberAdd(s.rest, data)

	if err != nil {
		return
	}

	// ToDo : Implement Member count?

	s.State().AddMember(ev.Data.GuildId, ev.Data)

	s.Publish(event.EventGuildMemberAdd, ev.Data)
}

type GuildMemberRemoveHandler struct{}

func (_ *GuildMemberRemoveHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildMemberRemove(s.rest, data)

	if err != nil {
		return
	}

	// ToDo : Implement Member count?

	s.State().RemoveMember(ev.Data.GuildId, ev.Data.User.Id)

	s.Publish(event.EventGuildMemberRemove, ev.Data)
}

type GuildMemberUpdateHandler struct{}

func (_ *GuildMemberUpdateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildMemberUpdate(s.rest, data)

	if err != nil {
		return
	}

	s.State().AddMember(ev.Data.GuildId, ev.Data)

	s.Publish(event.EventGuildMemberUpdate, ev.Data)
}

type GuildMembersChunkHandler struct{}

func (_ *GuildMembersChunkHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildMembersChunk(s.rest, data)

	if err != nil {
		return
	}

	for _, member := range ev.Data.Members {
		s.State().AddMember(ev.Data.GuildId, member)
	}

	s.Publish(event.EventGuildMembersChunk, ev.Data)
}

type GuildRoleCreateHandler struct{}

func (_ *GuildRoleCreateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildRoleCreate(s.rest, data)

	if err != nil {
		return
	}

	err = s.State().AddRole(ev.Data.GuildId, ev.Data.Role)

	if err != nil {
		return
	}

	s.Publish(event.EventGuildRoleCreate, ev.Data)
}

type GuildRoleUpdateHandler struct{}

func (_ *GuildRoleUpdateHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildRoleUpdate(s.rest, data)

	if err != nil {
		return
	}

	err = s.State().AddRole(ev.Data.GuildId, ev.Data.Role)

	if err != nil {
		return
	}

	s.Publish(event.EventGuildRoleUpdate, ev.Data)
}

type GuildRoleDeleteHandler struct{}

func (_ *GuildRoleDeleteHandler) Handle(s *Session, data []byte) {
	ev, err := event.NewGuildRoleDelete(s.rest, data)

	if err != nil {
		return
	}

	err = s.State().RemoveRole(ev.Data.GuildId, ev.Data.RoleId)

	if err != nil {
		return
	}

	s.Publish(event.EventGuildRoleDelete, ev.Data)
}
