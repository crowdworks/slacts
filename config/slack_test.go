package config

import (
	"os"
	"testing"
)

func TestNewSlackClientConfig(t *testing.T) {
	type envs struct {
		slackAPIToken string
	}

	cases := map[string]struct {
		envs envs
		want *SlackClientConfig
	}{
		"normal": {
			envs: envs{
				slackAPIToken: "slack_api_token",
			},
			want: &SlackClientConfig{
				token: "slack_api_token",
			},
		},
		"not set environment variable": {
			envs: envs{},
			want: &SlackClientConfig{},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			os.Setenv("SLACK_API_TOKEN", c.envs.slackAPIToken)
			defer os.Setenv("SLACK_API_TOKEN", "")

			got := NewSlackClientConfig()
			if got.Token() != c.want.Token() {
				t.Errorf("SlackClientConfig.Token() = %v, want %v", got.Token(), c.want.Token())
			}
		})
	}
}
