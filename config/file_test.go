package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadYaml(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]struct {
		file    string
		want    []TaskConfig
		wantErr bool
	}{
		"count": {
			file: filepath.Join(pwd, "./testdata/count.yml"),
			want: []TaskConfig{
				{
					Name:  "test_task",
					Kind:  "count",
					Query: "in:#general on:2018/12/03",
					Datadog: Datadog{
						Metric: "general.slack.count",
						Tags:   []string{"from:test", "env:test"},
					},
				},
			},
			wantErr: false,
		},
		"undefined kind": {
			file:    filepath.Join(pwd, "./testdata/no-kind.yml"),
			want:    nil,
			wantErr: true,
		},
		"not yaml": {
			file:    filepath.Join(pwd, "./testdata/unformatted.yml"),
			want:    nil,
			wantErr: true,
		},
		"invalid file path": {
			file:    filepath.Join(pwd, "./testdata/invalid-file-path"),
			want:    nil,
			wantErr: true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := ReadYaml(c.file)
			if (err != nil) != c.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if c.wantErr {
				return
			}

			if !reflect.DeepEqual(got, &c.want) {
				t.Errorf("Read() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestTaskConfig_DoesSendDatadog(t *testing.T) {
	type fields struct {
		Name    string
		Kind    string
		Query   string
		Datadog Datadog
	}
	cases := map[string]struct {
		fields fields
		want   bool
	}{
		"has datadog config": {
			fields: fields{
				Name:  "example",
				Kind:  "count",
				Query: "in:#general",
				Datadog: Datadog{
					Metric: "sample",
				},
			},
			want: true,
		},
		"empty datadog config": {
			fields: fields{
				Name:    "example",
				Kind:    "count",
				Query:   "in:#general",
				Datadog: Datadog{},
			},
			want: false,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			tc := &TaskConfig{
				Name:    c.fields.Name,
				Kind:    c.fields.Kind,
				Query:   c.fields.Query,
				Datadog: c.fields.Datadog,
			}
			if got := tc.DoesSendDatadog(); got != c.want {
				t.Errorf("TaskConfig.DoesSendDatadog() = %v, want %v", got, c.want)
			}
		})
	}
}
