package config

import (
	"fmt"
	"strings"
)

// Config contains all the configuration values
type Config struct {
	// PostgresURL the PostgresURL formatted like 'postgres://user:password@host:port/db'
	PostgresURL string

	// Channel is the name of the channel to intercept messages on
	Channel string
}

func (c *Config) Validate() error {
	var errs []string

	if c.PostgresURL == "" {
		errs = append(errs, "Missing PostgresURL")
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("Unable to start: %s", strings.Join(errs, ", "))
}
