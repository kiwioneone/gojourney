package rest

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
)

type VoiceHandler struct {
	rest *Client
}

func NewVoiceHandler(rest *Client) *VoiceHandler {
	return &VoiceHandler{rest: rest}
}

func (h *VoiceHandler) ListVoiceRegions() ([]*discord.VoiceRegion, error) {
	data, err := h.rest.Request(EndpointListVoiceRegions, "GET", nil, "application/json")

	if err != nil {
		return nil, err
	}

	var regions []*discord.VoiceRegion
	err = sonic.Unmarshal(data, &regions)

	if err != nil {
		return nil, err
	}

	return regions, nil
}
