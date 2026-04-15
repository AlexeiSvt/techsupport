package config

import "fmt"

func (c *Config) BuildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Storage.User, c.Storage.Password, c.Storage.Host, c.Storage.Port, c.Storage.Database, c.Storage.SSLMode)
}
