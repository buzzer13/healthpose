package icmp

import (
	"context"
	"fmt"
	probing "github.com/prometheus-community/pro-bing"
	"time"
)

// Config is the ICMP checker configuration settings container.
type Config struct {
	// Address is the IP address or hostname that is being checked.
	Address string
	// Count is the number of ICMP packets to send to the Address
	// Default is 1 packet.
	Count int
	// Interval is the wait time between each packet send.
	// Default is the value from probing.Pinger struct (1 second).
	Interval time.Duration
	// RequestTimeout is the duration that health check will try to consume published test message.
	// Default is the value from probing.Pinger struct (5 seconds).
	RequestTimeout time.Duration
}

// New creates new ICMP service health check that verifies the following:
// - domain name record can be resolved (if hostname is provided instead of the IP address)
// - ICMP echo request is sent and response is received
func New(config Config) func(ctx context.Context) error {
	if config.Count == 0 {
		config.Count = 1
	}

	return func(ctx context.Context) error {
		pinger, err := probing.NewPinger(config.Address)

		if err != nil {
			return fmt.Errorf("ping address resolve error: %w", err)
		}

		pinger.Count = config.Count
		pinger.Interval = config.Interval
		err = pinger.RunWithContext(ctx)

		if err != nil {
			return fmt.Errorf("ping error: %w", err)
		}

		return nil
	}
}
