package job

import (
	"context"
	"encoding/json"

	"github.com/Arriven/db1000n/logs"
)

// Args comment for linter
type Args = json.RawMessage

type Job = func(context.Context, *logs.Logger, Args) error

// Config comment for linter
type Config struct {
	Type  string
	Count int
	Args  Args
}

// Get job by type name
func Get(t string) (Job, bool) {
	res, ok := map[string]Job{
		"http":       httpJob,
		"tcp":        tcpJob,
		"udp":        udpJob,
		"syn-flood":  synFloodJob,
		"slow-loris": slowLorisJob,
		"packetgen":  packetgenJob,
	}[t]

	return res, ok
}

// BasicJobConfig comment for linter
type BasicJobConfig struct {
	IntervalMs int `json:"interval_ms,omitempty"`
	Count      int `json:"count,omitempty"`

	iter int
}

// Next comment for linter
func (c *BasicJobConfig) Next(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	if c.Count <= 0 {
		return true
	}

	c.iter++

	return c.iter <= c.Count
}

func panicHandler() {
	if err := recover(); err != nil {
		logs.Default.Warning("caught panic: %v\n some of the attacks may be unsupported on your system", err)
	}
}
