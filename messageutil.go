package gojourney

import (
	"regexp"

	"github.com/kiwioneone/gojourney/discord"
)

var re = regexp.MustCompile(`([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)

func MessageID(msg *discord.Message) string {
	if msg == nil {
		return ""
	}

	return msg.Id
}

func MessageAttachmentURL(msg *discord.Message) string {
	if msg == nil {
		return ""
	}

	if len(msg.Attachments) == 0 {
		return ""
	}
	return msg.Attachments[0].ProxyURL
}

func MessageJobID(msg *discord.Message) string {
	if msg == nil {
		return ""
	}

	if len(msg.Attachments) == 0 {
		return ""
	}

	//splitted := strings.Split(msg.Attachments[0].Filename, "_")

	return JobIDRegexp(msg.Attachments[0].Filename)
}

func JobIDRegexp(s string) string {
	match := re.FindStringSubmatch(s)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func IsMidjourneyReply(msg *discord.Message) bool {
	if msg.Author != nil && msg.Author.Id == MidjourneyApplicationID {
		return true
	}
	return false
}

func ReferenceMessageID(msg *discord.Message) string {
	if msg == nil {
		return ""
	}

	if msg.MessageReference == nil {
		return ""
	}

	return msg.MessageReference.MessageId
}

func ToCommandResult(msg *discord.Message) *CommandResult {
	if msg == nil {
		return nil
	}

	return &CommandResult{
		MessageID: MessageID(msg),
		ImageURL:  MessageAttachmentURL(msg),
		JobID:     MessageJobID(msg),
	}
}

func MessageEmbedDescription(msg *discord.Message) string {
	if msg == nil {
		return ""
	}

	if len(msg.Embeds) == 0 {
		return ""
	}

	return msg.Embeds[0].Description
}

func MessageCommandResult(msg *discord.Message) *CommandResult {
	if msg == nil {
		return nil
	}

	return &CommandResult{
		MessageID: MessageID(msg),
		ImageURL:  MessageAttachmentURL(msg),
		JobID:     MessageJobID(msg),
	}
}
