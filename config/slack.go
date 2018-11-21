package config

// SlackClientConfig is config object for slack client
type SlackClientConfig struct {
	// Slack API Token
	token string
}

// NewSlackClientConfig creates config object for slack client from environment variables
func NewSlackClientConfig() *SlackClientConfig {
	return &SlackClientConfig{
		token: getStringEnv("SLACK_API_TOKEN", ""),
	}
}

// Token returns slack api token
func (scc *SlackClientConfig) Token() string {
	return scc.token
}
