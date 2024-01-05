package gojourney

import (
	"strconv"

	"github.com/kiwioneone/gojourney/discord"
)

type Payload struct {
	Type          int         `json:"type"`
	ApplicationId string      `json:"application_id"`
	GuildId       string      `json:"guild_id"`
	ChannelId     string      `json:"channel_id"`
	SessionId     string      `json:"session_id"`
	Data          interface{} `json:"data"`
	MessageID     string      `json:"message_id,omitempty"`
	MessageFlags  int         `json:"message_flags,omitempty"`
	Nonce         string      `json:"nonce"`
}

type ApplicationCommandData struct {
	Version string                                             `json:"version"`
	Id      string                                             `json:"id"`
	Name    string                                             `json:"name"`
	Type    int                                                `json:"type"`
	Options []*discord.ApplicationCommandInteractionDataOption `json:"options"`
}

type ApplicationCommandIdVersion struct {
	Version string
	ID      string
}

func NewApplicationCommandPayload(guildID, channelID, sessionID string, commandName, commandParamsName, commandParamsValue string, commandVersionID ApplicationCommandIdVersion) *Payload {
	data := &ApplicationCommandData{
		Version: commandVersionID.Version,
		Id:      commandVersionID.ID,
		Name:    commandName,
		Type:    discord.ApplicationCommandChat,
	}

	if commandParamsName != "" {
		data.Options = []*discord.ApplicationCommandInteractionDataOption{
			{
				Name:  commandParamsName,
				Type:  discord.ApplicationCommandOptionString,
				Value: commandParamsValue,
			},
		}
	}

	return &Payload{
		Type:          discord.InteractionTypeApplicationCommand,
		ApplicationId: MidjourneyApplicationID,
		GuildId:       guildID,
		ChannelId:     channelID,
		SessionId:     sessionID,
		Nonce:         strconv.FormatInt(NextNonce(), 10),
		Data:          data,
	}
}

func NewMessageComponentPayload(guildID, channelID, sessionID string, messageID string, componentType discord.ComponentType, customID string) *Payload {
	return &Payload{
		Type:          discord.InteractionTypeMessageComponent,
		ApplicationId: MidjourneyApplicationID,
		GuildId:       guildID,
		ChannelId:     channelID,
		SessionId:     sessionID,
		Data: &discord.MessageComponentData{
			CustomId:      customID,
			ComponentType: componentType,
		},
		MessageID: messageID,
		Nonce:     strconv.FormatInt(NextNonce(), 10),
	}
}

func NewModalSubmitPayload(guildID, channelID, sessionID string, messageID string, customID string, components []discord.MessageComponent) *Payload {
	return &Payload{
		Type:          discord.InteractionTypeModalSubmit,
		ApplicationId: MidjourneyApplicationID,
		GuildId:       guildID,
		ChannelId:     channelID,
		SessionId:     sessionID,
		Data: &discord.ModalSubmitData{
			Id:         messageID,
			CustomId:   customID,
			Components: components,
		},
		Nonce: strconv.FormatInt(NextNonce(), 10),
	}
}
