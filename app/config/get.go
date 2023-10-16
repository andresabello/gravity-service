package config

func (c *Config) Get(key string) string {
	switch key {
	case "DBName":
		return c.DBName
	case "DBUser":
		return c.DBUser
	case "DBPassword":
		return c.DBPassword
	case "DBHost":
		return c.DBHost
	case "APIKey":
		return c.APIKey
	case "TwitterConsumerKey":
		return c.TwitterConsumerKey
	case "TwitterConsumerSecret":
		return c.TwitterConsumerSecret
	case "TwitterAccessToken":
		return c.TwitterAccessToken
	case "TwitterAccessSecret":
		return c.TwitterAccessSecret
	// Add other cases for additional configuration fields
	default:
		return ""
	}
}
