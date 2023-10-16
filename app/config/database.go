package config

import "fmt"

func (c *Config) ConstructDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBName,
	)
}
