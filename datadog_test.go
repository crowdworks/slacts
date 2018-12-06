package slacts_test

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/crowdworks/slacts"
	datadog "github.com/zorkian/go-datadog-api"
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

func TestNewDatadogClient(t *testing.T) {
	type args struct {
		apiKey     string
		appKey     string
		httpclient *http.Client
	}
	cases := map[string]struct {
		args args
		want *slacts.DatadogClient
	}{
		"default client": {
			args: args{
				apiKey:     "api_key",
				appKey:     "app_key",
				httpclient: nil,
			},
			want: &slacts.DatadogClient{
				Client: datadog.NewClient("api_key", "app_key"),
			},
		},
		"custom http client": {
			args: args{
				apiKey: "api_key",
				appKey: "app_key",
				httpclient: &http.Client{
					Transport:     nil,
					CheckRedirect: nil,
					Jar:           nil,
					Timeout:       10000 * time.Second,
				},
			},
			want: &slacts.DatadogClient{
				Client: datadogCustomClient("api_key", "app_key", &http.Client{
					Transport:     nil,
					CheckRedirect: nil,
					Jar:           nil,
					Timeout:       10000 * time.Second,
				}),
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := slacts.NewDatadogClient(c.args.apiKey, c.args.appKey, c.args.httpclient); !reflect.DeepEqual(got, c.want) {
				t.Errorf("NewDatadogClient() = %v, want %v", got, c.want)
			}
		})
	}
}

func datadogCustomClient(apiKey, appKey string, httpclient *http.Client) *datadog.Client {
	c := datadog.NewClient(apiKey, appKey)
	c.HttpClient = httpclient
	return c
}
