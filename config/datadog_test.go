package config

import (
	"os"
	"testing"
)

func TestNewDatadogConfig(t *testing.T) {
	type envs struct {
		appKey string
		apiKey string
	}

	cases := map[string]struct {
		envs       envs
		wantAPIKey string
		wantAppKey string
	}{
		"normal": {
			envs: envs{
				appKey: "datadog_app_token",
				apiKey: "datadog_api_token",
			},
			wantAppKey: "datadog_app_token",
			wantAPIKey: "datadog_api_token",
		},
		"not set env": {
			envs:       envs{},
			wantAPIKey: "",
			wantAppKey: "",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_ = os.Setenv("DATADOG_APP_KEY", c.envs.appKey)
			_ = os.Setenv("DATADOG_API_KEY", c.envs.apiKey)

			got := NewDatadogConfig()

			if got.AppKey() != c.wantAppKey {
				t.Errorf("AppKey() = %s, but want %s", got.AppKey(), c.wantAppKey)
			}

			if got.APIKey() != c.wantAPIKey {
				t.Errorf("APIKey() = %s, but want %s", got.APIKey(), c.wantAPIKey)
			}
		})
	}
}
