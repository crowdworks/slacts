package slacts_test

import (
	"errors"
	"testing"
	"time"

	"github.com/crowdworks/slacts"
)

type testDatadogClient struct {
	hasError bool
}

func (tdc *testDatadogClient) PostMetrics(metrics []slacts.DatadogMetric) error {
	if tdc.hasError {
		return errors.New("some error occurred")
	}

	return nil
}

func TestDatadogClient_PostMetrics(t *testing.T) {
	type fields struct {
		Client slacts.DatadogRequester
	}
	type args struct {
		metrics []slacts.DatadogMetric
	}
	cases := map[string]struct {
		fields  fields
		args    args
		wantErr bool
	}{
		"normal": {
			fields: fields{
				Client: &testDatadogClient{
					hasError: false,
				},
			},
			args: args{
				metrics: []slacts.DatadogMetric{
					{
						Metric: stringPointer("test.post.metric"),
						Points: []slacts.DatadogDataPoint{
							{
								float64Pointer(float64(time.Now().Unix())),
								float64Pointer(3.0),
							},
						},
						Tags: []string{"env:test", "channel:general"},
					},
				},
			},
			wantErr: false,
		},
		"nil metric": {
			fields: fields{
				Client: &testDatadogClient{
					hasError: false,
				},
			},
			args: args{
				metrics: nil,
			},
			wantErr: true,
		},
		"empty metric": {
			fields: fields{
				Client: &testDatadogClient{
					hasError: false,
				},
			},
			args: args{
				metrics: []slacts.DatadogMetric{},
			},
			wantErr: true,
		},
		"some server error": {
			fields: fields{
				Client: &testDatadogClient{
					hasError: true,
				},
			},
			args: args{
				metrics: []slacts.DatadogMetric{
					{
						Metric: stringPointer("test.post.metric"),
						Points: []slacts.DatadogDataPoint{
							{
								float64Pointer(float64(time.Now().Unix())),
								float64Pointer(3.0),
							},
						},
						Tags: []string{"env:test", "channel:general"},
					},
				},
			},
			wantErr: true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			dc := &slacts.DatadogClient{
				Client: c.fields.Client,
			}
			if err := dc.PostMetrics(c.args.metrics); (err != nil) != c.wantErr {
				t.Errorf("DatadogClient.PostMetrics() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}

func stringPointer(v string) *string {
	return &v
}

func float64Pointer(v float64) *float64 {
	return &v
}
