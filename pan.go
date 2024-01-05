package gojourney

import (
	"context"
	"fmt"
)

func (c *Client) PanDown(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::JOB::pan_down::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) PanUp(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::JOB::pan_up::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) PanLeft(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::JOB::pan_left::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) PanRight(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::JOB::pan_right::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}
