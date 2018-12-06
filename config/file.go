package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	kinds = []string{"count"}
)

// TaskConfig define what task should do
type TaskConfig struct {
	// Name of task
	Name string

	// Kind of task. only count is available
	Kind string

	// Query for slack message search
	Query string

	// Datadog config for metrics
	Datadog Datadog
}

// Datadog config for metrics
type Datadog struct {
	Metric string
	Tags   []string
}

// ReadYaml of given file path
func ReadYaml(file string) (*[]TaskConfig, error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var tasks []TaskConfig
	if err := viper.UnmarshalKey("tasks", &tasks); err != nil {
		return nil, err
	}

	for _, task := range tasks {
		if err := task.validate(); err != nil {
			return nil, err
		}
	}

	return &tasks, nil
}

func (tc *TaskConfig) validate() error {
	if err := tc.validateKind(); err != nil {
		return err
	}

	return nil
}

func (tc *TaskConfig) validateKind() error {
	for _, k := range kinds {
		if k == tc.Kind {
			return nil
		}
	}

	return fmt.Errorf("invalid kind. available: %v", kinds)
}

// DoesSendDatadog returns config has datadog config
func (tc *TaskConfig) DoesSendDatadog() bool {
	return tc.Datadog.Metric != ""
}
