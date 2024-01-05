package gojourney

import (
	"context"
	"fmt"
)

func (c *Client) Variation(ctx context.Context, messageID string, jobID string, index int) (result *CommandResult, err error) {
	customID := fmt.Sprintf("MJ::JOB::variation::%d::%s", index, jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) VaryStrong(ctx context.Context, messageID string, jobID string) (result *CommandResult, err error) {
	customID := fmt.Sprintf("MJ::JOB::high_variation::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) VarySubtle(ctx context.Context, messageID string, jobID string) (result *CommandResult, err error) {
	customID := fmt.Sprintf("MJ::JOB::low_variation::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}
