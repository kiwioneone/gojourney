package gojourney

import (
	"time"
)

// Epoch is the discord epoch in milliseconds.
const Epoch = 1420070400000

// NextNonce creates a new snowflake ID from the provided timestamp with worker id and sequence 0.
func NextNonce() int64 {
	timestamp := time.Now()
	return (timestamp.UnixMilli() - Epoch) << 22
}
