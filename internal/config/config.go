// Package config contains configuration structs and validation
package config

import (
	"fmt"
	"net/url"
	"strings"
)

// Config contains all the configuration values
type Config struct {
	// PostgresURL the PostgresURL formatted like 'postgres://user:password@host:port/db'
	PostgresURL string

	// Channel is the name of the channel to intercept messages on
	Channel string
}

// Validate makes sure all configuration is valid and the program can start
func (c *Config) Validate() error {
	var errs []string

	if c.PostgresURL == "" {
		errs = append(errs, "Missing PostgresURL")
	} else {
		_, err := url.Parse(c.PostgresURL)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}

	if c.Channel == "" {
		errs = append(errs, "Missing channel name")
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("Unable to start: %s", strings.Join(errs, ", "))
}
