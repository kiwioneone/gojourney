package gojourney

import (
	"context"
	"fmt"
)

func (c *Client) Upscale(ctx context.Context, messageID string, jobID string, index int) (result *CommandResult, err error) {
	customID := fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) Upscale2x(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::JOB::upsample_v5_2x::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) Upscale4x(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::JOB::upsample_v5_4x::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}
