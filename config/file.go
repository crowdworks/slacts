package config

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/spf13/viper"
)

var (
	// TimeNow can be define value from outside
	// If undefined, timeNow() generate real time now
	TimeNow *time.Time

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
	reader, err := parseTemplate(file)
	if err != nil {
		return nil, err
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(reader); err != nil {
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

func parseTemplate(file string) (io.Reader, error) {
	t, err := template.ParseFiles(file)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	type templateValues struct {
		// Yesterday format Y/m/d
		Yesterday string
	}

	r := &templateValues{
		Yesterday: timeNow().Add(-24 * time.Hour).Format("2006/01/02"),
	}
	if err := t.Execute(&b, r); err != nil {
		return nil, err
	}

	return &b, nil
}

// ReadYamlOption is type of ReadYaml option
type ReadYamlOption func(*Tasks) *Tasks

// OptionNameFilter is option for ReadYaml
// filter by task name
func OptionNameFilter(names []string) ReadYamlOption {
	return func(tasks *Tasks) *Tasks {
		var ts Tasks
		for _, task := range *tasks {
			for _, name := range names {
				if task.Name == name {
					ts = append(ts, task)
				}
			}
		}
		return &ts
	}
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

func timeNow() time.Time {
	if TimeNow != nil {
		return *TimeNow
	}

	return time.Now()
}
