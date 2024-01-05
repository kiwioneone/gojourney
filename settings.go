package gojourney

import (
	"context"
	"regexp"
	"time"

	"github.com/kiwioneone/gojourney/discord"
)

func NewSettingsCommandPayload(guildID, channelID, sessionID string) *Payload {
	return NewApplicationCommandPayload(
		guildID,
		channelID,
		sessionID,
		"settings",
		"",
		"",
		ApplicationCommandIdVersion{
			Version: "1166847114609958943",
			ID:      "1000850743479255081",
		})
}

type StylizeType string

var (
	StylizeLow      StylizeType = "low"
	StylizeMed      StylizeType = "med"
	StylizeHigh     StylizeType = "high"
	StylizeVeryHigh StylizeType = "very_high"
)

type VariationMode string

var (
	HighVariationMode VariationMode = "high"
	LowVariationMode  VariationMode = "low"
)

type SpeedMode string

var (
	TurboMode SpeedMode = "turbo"
	FastMode  SpeedMode = "fast"
	RelaxMode SpeedMode = "relax"
)

type Settings struct {
	DefaultModel  string        `json:"default_model,omitempty"`
	RawMode       bool          `json:"raw_mode,omitempty"`
	StylizeType   StylizeType   `json:"stylize_type,omitempty"`
	PublicMode    bool          `json:"public_mode,omitempty"`
	RemixMode     bool          `json:"remix_mode,omitempty"`
	VariationMode VariationMode `json:"variation_mode,omitempty"`
	StickyStyle   bool          `json:"sticky_style,omitempty"`
	SpeedMode     SpeedMode     `json:"speed_mode,omitempty"`
}

func (c *Client) Settings(ctx context.Context) (*Settings, error) {
	payload := NewSettingsCommandPayload(c.GuildID, c.ChannelID, c.ws.SessionID())
	messageFinder := &MessageFinder{Nonce: payload.Nonce, NoJob: true}

	ob := NewCommonMessageObserver(messageFinder.FilterMessage)
	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	err := c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return nil, err
	}

	msg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return nil, err
	}

	return getSettings(msg.Components), nil
}

func getSettings(components []discord.MessageComponent) *Settings {
	result := &Settings{}
	for _, component := range components {
		actions, ok := component.(*discord.ActionRows)
		if !ok || actions == nil {
			continue
		}
		for _, actionComponent := range actions.Components {
			menu, ok := actionComponent.(*discord.SelectMenu)
			if ok {
				ParseSettingMenu(menu, result)
			}

			button, ok := actionComponent.(*discord.Button)
			if ok {
				ParseSettingButton(button, result)
			}
		}
	}

	return result
}

func ParseSettingMenu(menu *discord.SelectMenu, settings *Settings) {
	if menu.CustomId != "MJ::Settings::VersionSelector" {
		return
	}

	for _, opt := range menu.Options {
		if opt.Default {
			settings.DefaultModel = ParseVersionModel(opt.Label)
		}
	}
}

func ParseSettingButton(button *discord.Button, settings *Settings) {
	if button == nil {
		return
	}

	isClicked := IsStyleClicked(button.Style)
	switch button.Label {
	case "RAW Mode":
		settings.RawMode = isClicked
	case "Stylize low":
		if isClicked {
			settings.StylizeType = StylizeLow
		}
	case "Stylize med":
		if isClicked {
			settings.StylizeType = StylizeMed
		}
	case "Stylize high":
		if isClicked {
			settings.StylizeType = StylizeHigh
		}
	case "Stylize very high":
		if isClicked {
			settings.StylizeType = StylizeVeryHigh
		}
	case "Public mode":
		settings.PublicMode = isClicked
	case "Remix mode":
		settings.RemixMode = isClicked
	case "High Variation Mode":
		if isClicked {
			settings.VariationMode = HighVariationMode
		}
	case "Low Variation Mode":
		if isClicked {
			settings.VariationMode = LowVariationMode
		}
	case "Sticky Style":
		settings.StickyStyle = isClicked
	case "Turbo mode":
		if isClicked {
			settings.SpeedMode = TurboMode
		}
	case "Fast mode":
		if isClicked {
			settings.SpeedMode = FastMode
		}
	case "Relax mode":
		if isClicked {
			settings.SpeedMode = RelaxMode
		}
	}
}

func IsStyleClicked(style discord.ButtonStyle) bool {
	return style == discord.ButtonStyleSuccess
}

func ParseVersionModel(str string) string {
	re := regexp.MustCompile(`\(([^)]+)\)`)
	result := re.FindStringSubmatch(str)
	if len(result) == 2 {
		return result[1]
	}
	return ""
}
