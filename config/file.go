package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	kinds = []string{"count"}
)

// Tasks is array of Task
type Tasks []Task

// Task define what task should do
type Task struct {
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
func ReadYaml(file string, opts ...ReadYamlOption) (*Tasks, error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	tasks := new(Tasks)
	if err := viper.UnmarshalKey("tasks", tasks); err != nil {
		return nil, err
	}

	for _, task := range *tasks {
		if err := task.validate(); err != nil {
			return nil, err
		}
	}

	for _, opt := range opts {
		tasks = opt(tasks)
	}

	return tasks, nil
}

// ReadYamlOption is type of ReadYaml option
type ReadYamlOption func(*Tasks) *Tasks

// OptionNameFilter is option for ReadYaml
// filter by task name
func OptionNameFilter(name string) ReadYamlOption {
	return func(tasks *Tasks) *Tasks {
		for i, task := range *tasks {
			if task.Name == name {
				tasks = unsetTask(tasks, i)
				return tasks
			}
		}
		return tasks
	}
}

// unsetTask remove given index of task
func unsetTask(tp *Tasks, i int) *Tasks {
	if i >= len(*tp) {
		return tp
	}

	t := *tp
	tasks := append(t[:i], t[i+1:]...)
	return &tasks
}

func (tc *Task) validate() error {
	if err := tc.validateKind(); err != nil {
		return err
	}

	return nil
}

func (tc *Task) validateKind() error {
	for _, k := range kinds {
		if k == tc.Kind {
			return nil
		}
	}

	return fmt.Errorf("invalid kind. available: %v", kinds)
}

// DoesSendDatadog returns config has datadog config
func (tc *Task) DoesSendDatadog() bool {
	return tc.Datadog.Metric != ""
}
