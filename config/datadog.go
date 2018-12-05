package config

// DatadogConfig is config object for datadog api client
type DatadogConfig struct {
	// application key
	appKey string

	// API key
	apiKey string
}

// NewDatadogConfig returns DatadogConfig object
func NewDatadogConfig() *DatadogConfig {
	return &DatadogConfig{
		appKey: getStringEnv("DATADOG_APP_KEY", ""),
		apiKey: getStringEnv("DATADOG_API_KEY", ""),
	}
}

// AppKey of Datadog
func (dc *DatadogConfig) AppKey() string {
	return dc.appKey
}

// APIKey of Datadog
func (dc *DatadogConfig) APIKey() string {
	return dc.apiKey
}
