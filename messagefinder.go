package gojourney

import (
	"errors"

	"github.com/kiwioneone/gojourney/discord"
)

var (
	stateTempCreate   = 0
	stateTempUpdate   = 1
	stateFinalMessage = 2
)

type MessageFinder struct {
	Nonce         string
	NoJob         bool
	jobID         string
	state         int
	tempMessageID string
}

// FilterMessage 过程参考 nonce->job_id->message_id https://bytedance.larkoffice.com/docx/GMK5dCck2oYLwJxcxSWckQJynUe
func (m *MessageFinder) FilterMessage(msg *discord.Message) (bool, error) {
	if msg == nil {
		return false, nil
	}

	if !IsMidjourneyReply(msg) {
		return false, nil
	}

	switch m.state {
	case stateTempCreate:
		if msg.Nonce == nil {
			return false, nil
		}

		msgNonce := msg.Nonce.(string)
		if msgNonce != m.Nonce {
			return false, nil
		}

		m.tempMessageID = MessageID(msg)
		m.state = stateTempUpdate

		embedDescription := MessageEmbedDescription(msg)
		if embedDescription != "" {
			return false, errors.New(embedDescription)
		}

		if len(msg.Components) == 0 {
			return false, nil
		}

		actions, ok := msg.Components[0].(*discord.ActionRows)
		if !ok {
			return false, nil
		}

		if len(actions.Components) == 0 {
			return false, nil
		}

		button, ok := actions.Components[0].(*discord.Button)
		if !ok {
			return false, nil
		}

		jobID := JobIDRegexp(button.CustomId)
		if jobID != "" {
			m.jobID = jobID
			m.state = stateFinalMessage
		}

	case stateTempUpdate:
		if MessageID(msg) != m.tempMessageID {
			return false, nil
		}

		if m.NoJob {
			return true, nil
		}

		jobID := MessageJobID(msg)
		if jobID != "" {
			m.jobID = jobID
			m.state = stateFinalMessage
			return false, nil
		}

		embedDescription := MessageEmbedDescription(msg)
		if embedDescription != "" {
			return false, errors.New(embedDescription)
		}

	case stateFinalMessage:
		if MessageID(msg) != m.tempMessageID {
			if MessageJobID(msg) == m.jobID {
				return true, nil
			}
		}
	default:
		return false, nil
	}
	return false, nil
}
