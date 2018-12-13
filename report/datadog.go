package report

import (
	"net/http"

	"github.com/pkg/errors"
	datadog "github.com/zorkian/go-datadog-api"
)

// DatadogMetric is alias of datadog.Metric
type DatadogMetric = datadog.Metric

// DatadogDataPoint is alias of datadog.DataPoint
type DatadogDataPoint = datadog.DataPoint

// DatadogRequester is interface of Datadog API client library
type DatadogRequester interface {
	PostMetrics([]DatadogMetric) error
}

// DatadogClient is client of Datadog API
type DatadogClient struct {
	Client DatadogRequester
}

// NewDatadogClient returns Datadog API client
func NewDatadogClient(apiKey, appKey string, httpclient *http.Client) *DatadogClient {
	client := datadog.NewClient(apiKey, appKey)
	if httpclient != nil {
		client.HttpClient = httpclient
	}

	return &DatadogClient{
		Client: client,
	}
}

// PostMetrics to datadog
func (dc *DatadogClient) PostMetrics(metrics []DatadogMetric) error {
	if len(metrics) == 0 {
		return errors.New("no metrics given")
	}

	if err := dc.Client.PostMetrics(metrics); err != nil {
		return err
	}

	return nil
}
