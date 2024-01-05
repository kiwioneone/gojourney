package gateway

// Intents is a bitfield of intents used to specify the gateway events you want to receive.
type Intents int

const (
	IntentGuilds Intents = 1 << iota
	IntentGuildMembers
	IntentGuildBans
	IntentGuildEmojisAndStickers
	IntentGuildIntegrations
	IntentGuildWebhooks
	IntentGuildInvites
	IntentGuildVoiceStates
	IntentGuildPresences
	IntentGuildMessages
	IntentGuildMessageReactions
	IntentGuildMessageTyping
	IntentDirectMessages
	IntentDirectMessageReactions
	IntentDirectMessageTyping
	IntentMessageContent
	IntentGuildScheduledEvents
	_
	_
	_
	IntentAutoModerationConfiguration
	IntentAutoModerationExecution

	IntentsGuild = IntentGuilds |
		IntentGuildMembers |
		IntentGuildBans |
		IntentGuildEmojisAndStickers |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildPresences |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentGuildScheduledEvents

	IntentsDirectMessage = IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping

	IntentsNonPrivileged = IntentGuilds |
		IntentGuildBans |
		IntentGuildEmojisAndStickers |
		IntentGuildIntegrations |
		IntentGuildWebhooks |
		IntentGuildInvites |
		IntentGuildVoiceStates |
		IntentGuildMessages |
		IntentGuildMessageReactions |
		IntentGuildMessageTyping |
		IntentDirectMessages |
		IntentDirectMessageReactions |
		IntentDirectMessageTyping |
		IntentGuildScheduledEvents |
		IntentAutoModerationConfiguration |
		IntentAutoModerationExecution

	IntentsPrivileged = IntentGuildMembers |
		IntentGuildPresences | IntentMessageContent

	IntentsAll = IntentsNonPrivileged |
		IntentsPrivileged

	IntentsDefault = IntentsNone

	IntentsNone int = 0
)
