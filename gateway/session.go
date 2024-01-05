package gateway

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/bytedance/sonic"

	ev "github.com/asaskevich/EventBus"
	"github.com/gorilla/websocket"

	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/gateway/event"
	"github.com/kiwioneone/gojourney/gateway/packet"
	"github.com/kiwioneone/gojourney/rest"
)

type Status int

const (
	StatusUnconnected Status = iota
	StatusConnecting
	StatusWaitingForHello
	StatusWaitingForReady
	StatusIdentifying
	StatusReady
	StatusResuming
	StatusDisconnected
)

type Session struct {
	sync.RWMutex

	options  *Options
	rest     *rest.Client
	presence *packet.PresenceUpdate
	user     *discord.User
	bus      *ev.EventBus
	state    *State
	status   Status
	Error    error

	// ws conn
	connMu   sync.Mutex
	conn     *websocket.Conn
	handlers map[event.EventType]EventHandler

	// Discord gateway fields
	sessionID         string
	heartbeatTicker   *time.Ticker
	heartbeatInterval time.Duration
	lastHeartbeatAck  time.Time
	lastHeartbeatSent time.Time
	lastSequence      int64

	// voice conn
	VoiceConnections map[string]*VoiceConnection

	// Rest handlers
	Application *rest.ApplicationHandler
	Channel     *rest.ChannelHandler
	Emoji       *rest.EmojiHandler
	Guild       *rest.GuildHandler
	Interaction *rest.InteractionHandler
	Invite      *rest.InviteHandler
	Template    *rest.TemplateHandler
	User        *rest.UserHandler
	Voice       *rest.VoiceHandler
	Webhook     *rest.WebhookHandler
}

func NewSession(options *Options) *Session {
	s := new(Session)

	s.options = options
	s.presence = packet.NewPresenceUpdate(nil, discord.StatusTypeOnline)
	s.user = new(discord.User)
	s.rest = rest.NewClient(options.Token)
	s.bus = ev.New().(*ev.EventBus)
	s.state = NewState(s)
	s.status = StatusUnconnected

	// voice conn
	s.VoiceConnections = make(map[string]*VoiceConnection)

	s.Application = rest.NewApplicationHandler(s.rest)
	s.Channel = rest.NewChannelHandler(s.rest)
	s.Emoji = rest.NewEmojiHandler(s.rest)
	s.Guild = rest.NewGuildHandler(s.rest)
	s.Interaction = rest.NewInteractionHandler(s.rest)
	s.Invite = rest.NewInviteHandler(s.rest)
	s.Template = rest.NewTemplateHandler(s.rest)
	s.User = rest.NewUserHandler(s.rest)
	s.Voice = rest.NewVoiceHandler(s.rest)
	s.Webhook = rest.NewWebhookHandler(s.rest)

	s.registerHandlers()

	return s
}

func (s *Session) registerHandlers() {
	s.handlers = map[event.EventType]EventHandler{
		event.EventReady:   &ReadyHandler{},
		event.EventResumed: &ResumedHandler{},
		// Application events
		event.EventApplicationCommandPermissionsUpdate: &ApplicationCommandPermissionsUpdateHandler{},
		// AutoModeration events
		event.EventAutoModerationRuleCreate:      &AutoModerationRuleCreateHandler{},
		event.EventAutoModerationRuleDelete:      &AutoModerationRuleDeleteHandler{},
		event.EventAutoModerationRuleUpdate:      &AutoModerationRuleUpdateHandler{},
		event.EventAutoModerationActionExecution: &AutoModerationActionExecutionHandler{},
		event.EventChannelCreate:                 &ChannelCreateHandler{},
		event.EventChannelUpdate:                 &ChannelUpdateHandler{},
		event.EventChannelDelete:                 &ChannelDeleteHandler{},
		event.EventChannelPinsUpdate:             &ChannelPinsUpdateHandler{},
		event.EventThreadCreate:                  &ThreadCreateHandler{},
		event.EventThreadUpdate:                  &ThreadUpdateHandler{},
		event.EventThreadDelete:                  &ThreadDeleteHandler{},
		event.EventThreadListSync:                &ThreadListSyncHandler{},
		event.EventThreadMemberUpdate:            &ThreadMemberUpdateHandler{},
		event.EventThreadMembersUpdate:           &ThreadMembersUpdateHandler{},
		event.EventGuildStickersUpdate:           &GuildStickersUpdateHandler{},
		event.EventGuildIntegrationsUpdate:       &GuildIntegrationsUpdateHandler{},
		event.EventGuildMemberAdd:                &GuildMemberAddHandler{},
		event.EventGuildMemberRemove:             &GuildMemberRemoveHandler{},
		event.EventGuildMemberUpdate:             &GuildMemberUpdateHandler{},
		event.EventGuildMembersChunk:             &GuildMembersChunkHandler{},
		event.EventGuildRoleCreate:               &GuildRoleCreateHandler{},
		event.EventGuildRoleUpdate:               &GuildRoleUpdateHandler{},
		event.EventGuildRoleDelete:               &GuildRoleDeleteHandler{},

		event.EventGuildCreate:       &GuildCreateHandler{},
		event.EventGuildUpdate:       &GuildUpdateHandler{},
		event.EventGuildDelete:       &GuildDeleteHandler{},
		event.EventGuildBanAdd:       &GuildBanAddHandler{},
		event.EventGuildBanRemove:    &GuildBanRemoveHandler{},
		event.EventGuildEmojisUpdate: &GuildEmojisUpdateHandler{},
		event.EventMessageCreate:     &MessageCreateHandler{},
		event.EventMessageUpdate:     &MessageUpdateHandler{},
		event.EventMessageDelete:     &MessageDeleteHandler{},
		event.EventPresenceUpdate:    &PresenceUpdateHandler{},
		event.EventInteractionCreate: &InteractionCreateHandler{},
		event.EventVoiceStateUpdate:  &VoiceStateUpdateHandler{},
		event.EventVoiceServerUpdate: &VoiceServerUpdateHandler{},
	}
}

// JoinVoiceChannel joins a voice channel.
func (s *Session) JoinVoiceChannel(guildId, channelId string, muted, deafened bool) (*VoiceConnection, error) {
	s.RLock()
	vConn, ok := s.VoiceConnections[guildId]
	s.RUnlock()

	if !ok {
		vConn = &VoiceConnection{}

		s.Lock()
		s.VoiceConnections[guildId] = vConn
		s.Unlock()
	}

	vConn.Lock()
	vConn.GuildId = guildId
	vConn.ChannelId = channelId
	vConn.deaf = deafened
	vConn.mute = muted
	vConn.session = s
	vConn.Unlock()

	voiceStateUpdate := packet.NewVoiceStateUpdate(guildId, channelId, muted, deafened)
	if err := s.Send(voiceStateUpdate); err != nil {
		return nil, err
	}

	if err := vConn.wait(); err != nil {
		vConn.Close()
		return nil, err
	}

	return vConn, nil
}

func (s *Session) JoinVoiceChannelIncomplete(guildId, channelId string, muted, deafened bool) error {
	return s.Send(packet.NewVoiceStateUpdate(guildId, channelId, muted, deafened))
}

// Login connects the session to the gateway.
func (s *Session) Login() error {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.conn != nil {
		return errors.New("session is already connected")
	}

	s.status = StatusConnecting
	s.lastHeartbeatSent = time.Now().UTC()

	conn, rs, err := websocket.DefaultDialer.Dial(rest.GatewayUrl, nil)
	if err != nil {
		body := "null"

		if rs != nil && rs.Body != nil {
			defer func() {
				_ = rs.Body.Close()
			}()

			rawBody, bErr := io.ReadAll(rs.Body)
			if bErr != nil {
				return err
			}

			body = string(rawBody)
		}

		s.Error = fmt.Errorf("error while connecting to the gateway : %s", body)
		return s.Error
	}

	conn.SetCloseHandler(func(code int, text string) error {
		closeCode := packet.CloseEventCode(code)

		if !closeCode.ShouldReconnect() {

			s.Error = &websocket.CloseError{
				Code: int(closeCode),
				Text: text,
			}
			return s.Error
			//return fmt.Errorf("error connecting to gateway : %d %s", code, text)
		}

		return nil
	})

	s.conn = conn
	s.status = StatusWaitingForHello

	go s.listen(conn)

	return nil
}

func (s *Session) listen(conn *websocket.Conn) {
loop:
	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			s.connMu.Lock()
			sameConnection := s.conn == conn
			s.connMu.Unlock()

			if !sameConnection {
				return
			}

			reconnect := true

			if closeError, ok := err.(*websocket.CloseError); ok {
				closeCode := packet.CloseEventCode(closeError.Code)
				reconnect = closeCode.ShouldReconnect()
			} else if errors.Is(err, net.ErrClosed) {
				reconnect = false
			}

			s.CloseWithCode(websocket.CloseServiceRestart, "reconnecting")
			if reconnect {
				go s.reconnect()

				break loop
			}
		}

		pk, err := packet.NewPacket(msg)

		if err != nil {
			return
		}

		opcode, e := pk.Opcode, pk.Event

		switch opcode {
		case packet.OpHello:
			s.connMu.Lock()
			s.lastHeartbeatAck = time.Now().UTC()
			s.connMu.Unlock()

			hello, err := packet.NewHello(msg)

			if err != nil {
				return
			}

			go s.startHeartbeat()

			s.connMu.Lock()
			s.heartbeatInterval = time.Duration(hello.Data.HeartbeatInterval) * time.Millisecond
			lastSequence := s.lastSequence
			sessionID := s.sessionID

			token := s.options.Token
			intents := s.options.Intents
			s.connMu.Unlock()

			if lastSequence == 0 || sessionID == "" {
				s.connMu.Lock()
				s.status = StatusIdentifying
				s.connMu.Unlock()

				identify := packet.NewIdentify(token, int(intents))

				if err = s.Send(identify); err != nil {
					return
				}

				s.connMu.Lock()
				s.status = StatusWaitingForReady
				s.connMu.Unlock()
			} else {
				resume := packet.NewResume(token, sessionID, lastSequence)

				if err = s.Send(resume); err != nil {
					return
				}
			}

		case packet.OpDispatch:
			s.connMu.Lock()
			s.lastSequence = pk.Sequence
			s.connMu.Unlock()

			if e != "" {
				s.connMu.Lock()
				s.lastSequence = pk.Sequence
				handler, exists := s.handlers[event.EventType(e)]
				s.connMu.Unlock()

				if exists {
					go handler.Handle(s, msg)
				} else {
					fmt.Println("Unhandled event : " + e)
				}
			}

		case packet.OpHeartbeat:
			s.sendHeartbeat()

		case packet.OpReconnect:
			s.CloseWithCode(websocket.CloseServiceRestart, "reconnecting")
			go s.reconnect()

			break loop

		case packet.OpInvalidSession:
			var shouldResume = false

			err = sonic.Unmarshal(pk.Data, &shouldResume)
			if err != nil {
				shouldResume = false
			}

			code := websocket.CloseNormalClosure
			if shouldResume {
				code = websocket.CloseServiceRestart
			} else {
				s.connMu.Lock()
				s.sessionID = ""
				s.lastSequence = 0
				s.connMu.Unlock()
			}

			s.CloseWithCode(code, "invalid session")

			go s.reconnect()

			break loop

		case packet.OpHeartbeatAck:
			s.connMu.Lock()
			s.lastHeartbeatAck = time.Now().UTC()
			s.connMu.Unlock()
		}
	}
}

func (s *Session) startHeartbeat() {
	s.connMu.Lock()
	heartbeatTicker := time.NewTicker(s.heartbeatInterval)
	s.heartbeatTicker = heartbeatTicker
	s.connMu.Unlock()

	defer heartbeatTicker.Stop()

	for range heartbeatTicker.C {
		s.sendHeartbeat()
	}
}

func (s *Session) sendHeartbeat() {
	s.connMu.Lock()
	lastSequence := s.lastSequence
	s.connMu.Unlock()

	heartbeat := packet.NewHeartbeat(lastSequence)

	if err := s.Send(heartbeat); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return
		}

		s.CloseWithCode(websocket.CloseServiceRestart, "heartbeat timeout")

		go s.reconnect()

		return
	}

	s.connMu.Lock()
	s.lastHeartbeatSent = time.Now().UTC()
	s.connMu.Unlock()
}

func (s *Session) reconnect() {
	wait := time.Duration(5)

	for {
		fmt.Println("Reconnecting")

		err := s.Login()

		if err == nil {
			fmt.Println("Reconnected")

			s.RLock()
			voiceConnections := s.VoiceConnections
			s.RUnlock()

			for _, v := range voiceConnections {
				v.reconnect()
			}

			return
		}

		<-time.After(wait)

		wait *= 2

		if wait > 300 {
			wait = 300
		}
	}
}

// Send sends a packet to the gateway.
func (s *Session) Send(v interface{}) error {
	s.connMu.Lock()
	defer s.connMu.Unlock()

	return s.conn.WriteJSON(v)
}

// SetActivity sets the activity of the session.
func (s *Session) SetActivity(activity *discord.Activity) error {
	s.Lock()
	s.presence.Data.Activities[0] = activity
	s.Unlock()

	s.RLock()
	defer s.RUnlock()

	return s.Send(s.presence)
}

// SetStatus sets the status of the session.
func (s *Session) SetStatus(status discord.StatusType) error {
	s.Lock()
	s.presence.Data.Status = status
	s.Unlock()

	s.RLock()
	defer s.RUnlock()

	return s.Send(s.presence)
}

// UpdatePresence updates the status and activity of the session.
func (s *Session) UpdatePresence(status *packet.PresenceUpdate) error {
	s.Lock()
	s.presence = status
	s.Unlock()

	return s.Send(status)
}

// Latency returns the latency of the session.
func (s *Session) Latency() time.Duration {
	s.connMu.Lock()
	lastHeartbeatAck := s.lastHeartbeatAck
	lastHeartbeatSent := s.lastHeartbeatSent
	s.connMu.Unlock()

	return lastHeartbeatAck.Sub(lastHeartbeatSent)
}

// Close closes the session.
func (s *Session) Close() {
	s.CloseWithCode(websocket.CloseNormalClosure, "Shutting down")
}

// CloseWithCode closes the session with a specific close code and message.
func (s *Session) CloseWithCode(code int, message string) {
	s.connMu.Lock()
	heartbeatTicker := s.heartbeatTicker
	s.connMu.Unlock()

	if heartbeatTicker != nil {
		heartbeatTicker.Stop()
		heartbeatTicker = nil
	}

	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.conn != nil {
		_ = s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message))

		_ = s.conn.Close()

		s.conn = nil

		if code == websocket.CloseNormalClosure || code == websocket.CloseGoingAway {
			s.sessionID = ""
			s.lastSequence = 0
		}
	}
}

// Bus returns the event bus of the session.
func (s *Session) Bus() *ev.EventBus {
	s.RLock()
	defer s.RUnlock()

	return s.bus
}

// Me returns the current user of the session.
func (s *Session) Me() *discord.User {
	s.RLock()
	defer s.RUnlock()

	return s.user
}

// State returns the current state of the session.
func (s *Session) State() *State {
	s.RLock()
	defer s.RUnlock()

	return s.state
}

// Status returns the current status of the session.
func (s *Session) Status() Status {
	s.RLock()
	defer s.RUnlock()

	return s.status
}

// On registers a callback for an event type.
func (s *Session) On(ev event.EventType, fn any) error {
	return s.Bus().SubscribeAsync(ev.String(), fn, false)
}

// Publish publishes an event to the event bus.
func (s *Session) Publish(ev event.EventType, args ...any) {
	s.Bus().Publish(ev.String(), args...)
}

func (s *Session) SessionID() string {
	return s.sessionID
}
