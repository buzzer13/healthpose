package dns

import (
	"context"
	"fmt"
	"net"
	"time"
)

type RecordType string

const (
	RecordTypeA     RecordType = "a"
	RecordTypeCNAME RecordType = "cname"
	RecordTypePTR   RecordType = "ptr"
	RecordTypeTXT   RecordType = "txt"
)

// Config is the DNS checker configuration settings container.
type Config struct {
	// Address is the IP address or hostname that is being checked.
	Address string
	// Server is the server that resolves hostname.
	Server string
	// Type is the DNS record type
	// Default is A record.
	Type RecordType
	// RequestTimeout is the duration that health check will try to consume published test message.
	// Default is 5 seconds.
	RequestTimeout time.Duration
	// FallbackDelay
	// Default is the value from net.Dialer struct (300ms).
	FallbackDelay time.Duration
}

const defaultRequestTimeout = 300 * time.Millisecond
const defaultRecordType = RecordTypeA

// New creates new DNS service health check that verifies the following:
// - domain name record can be resolved
func New(config Config) func(ctx context.Context) error {
	if config.RequestTimeout == 0 {
		config.RequestTimeout = defaultRequestTimeout
	}

	if config.Type == "" {
		config.Type = defaultRecordType
	}

	switch config.Type {
	case RecordTypeA:
		return Resolver(config, func(res *net.Resolver, ctx context.Context) error {
			_, err := res.LookupHost(ctx, config.Address)
			return err
		})
	case RecordTypeCNAME:
		return Resolver(config, func(res *net.Resolver, ctx context.Context) error {
			_, err := res.LookupCNAME(ctx, config.Address)
			return err
		})
	case RecordTypePTR:
		return Resolver(config, func(res *net.Resolver, ctx context.Context) error {
			_, err := res.LookupAddr(ctx, config.Address)
			return err
		})
	case RecordTypeTXT:
		return Resolver(config, func(res *net.Resolver, ctx context.Context) error {
			_, err := res.LookupTXT(ctx, config.Address)
			return err
		})
	default:
		return Resolver(config, func(res *net.Resolver, ctx context.Context) error {
			return fmt.Errorf("invalid domain record type: %s", config.Type)
		})
	}
}

type CheckerFunc func(res *net.Resolver, ctx context.Context) error

func Resolver(config Config, checker CheckerFunc) func(ctx context.Context) error {
	res := net.Resolver{
		PreferGo: true,
	}

	if config.Server != "" {
		res.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout:       config.RequestTimeout,
				FallbackDelay: config.FallbackDelay,
			}

			return d.DialContext(ctx, network, config.Server)
		}
	}

	return func(ctx context.Context) error {
		err := checker(&res, ctx)

		if err != nil {
			return fmt.Errorf("domain health check failed: %w", err)
		}

		return nil
	}
}
