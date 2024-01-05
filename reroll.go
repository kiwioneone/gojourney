package gojourney

import (
	"context"
	"fmt"
)

func (c *Client) Reroll(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::JOB::reroll::0::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}
