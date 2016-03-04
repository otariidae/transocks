package transocks

import (
	"errors"
	"fmt"
	"net"
	"net/url"
)

const (
	// NAT mode
	ModeNAT = "nat"
)

// Config keeps configurations for Server.
type Config struct {
	// Listen is the listening address.
	// e.g. "localhost:1081"
	Listen string

	// ProxyURL is the URL for upstream proxy.
	//
	// For SOCKS5, URL looks like "socks5://USER:PASSWORD@HOST:PORT".
	//
	// For HTTP proxy, URL looks like "http://USER:PASSWORD@HOST:PORT".
	// The HTTP proxy must support CONNECT method.
	ProxyURL *url.URL

	// Mode determines how clients are routed to transocks.
	// Default is "nat".  No other options are available at this point.
	Mode string

	// Dialer is the base dialer to connect to the proxy server.
	// The server uses the default dialer if this is nil.
	Dialer *net.Dialer
}

// NewConfig creates and initializes a new Config.
func NewConfig() *Config {
	c := new(Config)
	c.Mode = ModeNAT
	return c
}

// validate validates the configuration.
// It returns non-nil error if the configuration is not valid.
func (c *Config) validate() error {
	if len(c.Listen) == 0 {
		return errors.New("Listen is empty")
	}
	if c.ProxyURL == nil {
		return errors.New("ProxyURL is nil")
	}
	if c.Mode != ModeNAT {
		return fmt.Errorf("Unknown mode: %s", c.Mode)
	}
	return nil
}
