package discord

import (
	"time"

	"github.com/bytedance/sonic"
)

type MessageActivityType int

const (
	MessageActivityTypeJoin MessageActivityType = iota + 1
	MessageActivityTypeSpectate
	MessageActivityTypeListen
	MessageActivityTypeJoinRequest
)

type MessageFlag int

const (
	MessageFlagCrossposted MessageFlag = 1 << iota
	MessageFlagIsCrosspost
	MessageFlagSuppressEmbeds
	MessageFlagSourceMessageDeleted
	MessageFlagUrgent
	MessageFlagHasThread
	MessageFlagEphemeral
	MessageFlagLoading
	MessageFlagFailedToMentionSomeRolesInThreads
)

type AllowedMentionsType string

const (
	AllowedMentionsRoleMentions AllowedMentionsType = "roles"
	AllowedMentionsUserMentions AllowedMentionsType = "users"
	AllowedMentionsEveryone     AllowedMentionsType = "everyone"
)

type MessageType int

const (
	MessageTypeDefault MessageType = iota
	MessageTypeRecipientAdd
	MessageTypeRecipientRemove
	MessageTypeCall
	MessageTypeChannelNameChange
	MessageTypeChannelIconChange
	MessageTypeChannelPinnedMessage
	MessageTypeUserJoin
	MessageTypeGuildBoost
	MessageTypeGuildBoostTier1
	MessageTypeGuildBoostTier2
	MessageTypeGuildBoostTier3
	MessageTypeChannelFollowAdd
	MessageTypeGuildDiscoveryAdd
	MessageTypeGuildDiscoveryDisqualified
	MessageTypeGuildDiscoveryRequalified
	MessageTypeGuildDiscoveryGracePeriodInitialWarning
	MessageTypeGuildDiscoveryGracePeriodFinalWarning
	MessageTypeThreadCreated
	MessageTypeReply
	MessageTypeChatInputCommand
	MessageTypeThreadStarterMessage
	MessageTypeGuildInviteReminder
	MessageTypeContextMenuCommand
	MessageTypeAutoModerationAction
	MessageTypeRoleSubscription
	MessageTypeInteractionPremiumUpsell
	MessageTypeGuildApplicationPremiumSubscription = 32
)

type ChannelType int

const (
	ChannelTypeText ChannelType = iota
	ChannelTypeDM
	ChannelTypeVoice
	ChannelTypeGroupDM
	ChannelTypeCategory
	ChannelTypeNews
	ChannelTypeNewsThread
	ChannelTypePublicThread
	ChannelTypePrivateThread
	ChannelTypeStageVoice
	ChannelTypeDirectory
	ChannelTypeForum // Still in development
)

type VideoQuality int

const (
	_                VideoQuality = iota
	VideoQualityAuto              // Discord chooses the quality for optimal performance
	VideoQualityFull              // 720p, 1080p, or higher quality
)

type ChannelFlags int

const (
	ChannelFlagsPinned ChannelFlags = 1 << 1
)

type FollowedChannel struct {
	ChannelId string `json:"channel_id"`
	WebhookId string `json:"webhook_id"`
}

type Overwrite struct {
	Id    string                `json:"id"`           // role or user id
	Type  int                   `json:"type"`         // 0 for role, 1 for user
	Allow BitwisePermissionFlag `json:"allow,string"` // permission bit set
	Deny  BitwisePermissionFlag `json:"deny,string"`  // permission bit set
}

type MessageActivity struct {
	Type    string `json:"type"`
	PartyId string `json:"party_id"`
}

type MessageReference struct {
	MessageId      string `json:"message_id,omitempty"`
	ChannelId      string `json:"channel_id,omitempty"`
	GuildId        string `json:"guild_id,omitempty"`
	FailIfNotExist bool   `json:"fail_if_not_exist,omitempty"`
}

type Reaction struct {
	Count int    `json:"count"`
	Me    bool   `json:"me"`
	Emoji *Emoji `json:"emoji"`
}

type ThreadMetadata struct {
	Archived            bool       `json:"archived"`
	AutoArchiveDuration int        `json:"auto_archive_duration"`
	ArchiveTimestamp    *time.Time `json:"archive_timestamp"`
	Locked              bool       `json:"locked"`
	Invitable           bool       `json:"invitable,omitempty"`
	CreateTimestamp     *time.Time `json:"create_timestamp,omitempty"`
}

type ThreadMember struct {
	Id            string     `json:"id,omitempty"`
	UserId        string     `json:"user_id,omitempty"`
	JoinTimestamp *time.Time `json:"join_timestamp"`
	Flags         int        `json:"flags"`
	GuildId       string     `json:"guild_id,omitempty"`
}

type Attachment struct {
	Id       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	Data     []byte `json:"-"`
	ProxyURL string `json:"proxy_url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type ChannelMention struct {
	Id      string      `json:"id"`
	GuildId string      `json:"guild_id"`
	Type    ChannelType `json:"type"`
	Name    string      `json:"name"`
}

type AllowedMentions struct {
	Parse        []*AllowedMentions `json:"parse"`
	Roles        []string           `json:"roles"`
	Users        []string           `json:"users"`
	RepliedUsers bool               `json:"replied_users"`
}

type Message struct {
	Id                string              `json:"id"`
	ChannelId         string              `json:"channel_id"`
	GuildId           string              `json:"guild_id,omitempty"`
	Author            *User               `json:"author"`
	Member            *GuildMember        `json:"member"`
	Content           string              `json:"content"`
	Timestamp         *time.Time          `json:"timestamp"`
	EditedTimestamp   *time.Time          `json:"edited_timestamp"`
	Tts               bool                `json:"tts"`
	MentionEveryone   bool                `json:"mention_everyone"`
	Mentions          []*User             `json:"mentions"`
	MentionRoles      []string            `json:"mention_roles"`
	MentionChannels   []*Channel          `json:"mention_channels,omitempty"`
	Attachments       []*Attachment       `json:"attachments"`
	Embeds            []*Embed            `json:"embeds"`
	Reactions         []*Reaction         `json:"reactions"`
	Nonce             interface{}         `json:"nonce,omitempty"` // integer or string
	Pinned            bool                `json:"pinned"`
	WebhookId         string              `json:"webhook_id,omitempty"`
	Type              MessageType         `json:"type"`
	Activity          *MessageActivity    `json:"activity,omitempty"`
	Application       *Application        `json:"application,omitempty"`
	ApplicationId     string              `json:"application_id,omitempty"`
	MessageReference  *MessageReference   `json:"message_reference,omitempty"`
	Flags             MessageFlag         `json:"flags,omitempty"`
	ReferencedMessage *Message            `json:"referenced_message,omitempty"`
	Interaction       *MessageInteraction `json:"interaction,omitempty"`
	Thread            *Channel            `json:"thread,omitempty"`
	Components        []MessageComponent  `json:"components"`
	StickerItems      []*StickerItem      `json:"sticker_items,omitempty"`
	Position          int                 `json:"position,omitempty"`
}

type rawMessage Message

// UnmarshalJSON ...
func (m *Message) UnmarshalJSON(data []byte) error {
	var v struct {
		rawMessage
		Components []unmarshalableMessageComponent `json:"components"`
	}

	err := sonic.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*m = Message(v.rawMessage)

	m.Components = make([]MessageComponent, len(v.Components))

	for i, v := range v.Components {
		m.Components[i] = v.MessageComponent
	}

	return err
}

type Channel struct {
	Id                         string                `json:"id"`
	Type                       ChannelType           `json:"type"`
	GuildId                    string                `json:"guild_id,omitempty"`
	Position                   int                   `json:"position,omitempty"`
	PermissionOverwrites       []*Overwrite          `json:"permission_overwrites,omitempty"`
	Name                       string                `json:"name,omitempty"`
	Topic                      string                `json:"topic,omitempty"`
	Nsfw                       bool                  `json:"nsfw,omitempty"`
	LastMessageId              string                `json:"last_message_id,omitempty"`
	Bitrate                    int                   `json:"bitrate,omitempty"`
	UserLimit                  int                   `json:"user_limit,omitempty"`
	RateLimitPerUser           int                   `json:"rate_limit_per_user,omitempty"`
	Recipients                 []*User               `json:"recipients,omitempty"`
	Icon                       string                `json:"icon,omitempty"`
	OwnerId                    string                `json:"owner_id,omitempty"`
	ApplicationId              string                `json:"application_id,omitempty"`
	ParentId                   string                `json:"parent_id,omitempty"`
	LastPinTimestamp           *time.Time            `json:"last_pin_timestamp,omitempty"`
	RtcRegion                  string                `json:"rtc_region,omitempty"` // 	voice region id for the voice channel, automatic when set to null
	VideoQualityMode           VideoQuality          `json:"video_quality_mode,omitempty"`
	MessageCount               int                   `json:"message_count,omitempty"`
	MemberCount                int                   `json:"member_count,omitempty"`
	ThreadMetadata             *ThreadMetadata       `json:"thread_metadata,omitempty"`
	Member                     *ThreadMember         `json:"member,omitempty"`
	DefaultAutoArchiveDuration int                   `json:"default_auto_archive_duration,omitempty"`
	Permissions                BitwisePermissionFlag `json:"permissions,string,omitempty"`
	Flags                      ChannelFlags          `json:"flags,omitempty"`
	TotalMessageSent           int                   `json:"total_message_sent,omitempty"`
}

// ChannelPinsUpdateEventFields is used by the CHANNEL_PINS_UPDATE event
type ChannelPinsUpdateEventFields struct {
	GuildId          string     `json:"guild_id,omitempty"`
	ChannelId        string     `json:"channel_id"`
	LastPinTimestamp *time.Time `json:"last_pin_timestamp,omitempty"`
}

// ThreadListSyncEventFields is used by the THREAD_LIST_SYNC event
type ThreadListSyncEventFields struct {
	GuildId    string   `json:"guild_id"`
	ChannelIds []string `json:"channel_ids,omitempty"`
	Threads    []*Channel
	Members    []*ThreadMember
}

// ThreadMembersUpdateEventFields is used by the THREAD_MEMBERS_UPDATE event
type ThreadMembersUpdateEventFields struct {
	Id                string          `json:"id"`
	GuildId           string          `json:"guild_id"`
	MemberCount       int             `json:"member_count"`
	AddedMembers      []*ThreadMember `json:"added_members,omitempty"`
	RemovedMembersIds []string        `json:"removed_member_ids,omitempty"`
}
