package gojourney

import "errors"

const MidjourneyApplicationID = "936929561302675456"
const DiscordInteractionURL = "https://discord.com/api/v9/interactions"

var ErrorNotLogin = errors.New("not login")
var ErrorNotFoundJobID = errors.New("no job id")
var ErrorFactorInvalid = errors.New("factor need be (1,2]")
var ErrorEmptyEmbeds = errors.New("empty embeds")
var ErrorInvalidAssistBot = errors.New("no invalid assist")
