package slacts

import (
	"context"
	"net/http"

	"github.com/nlopes/slack"
)

// SlackRequester is interface that wrap slack.Client only methods this package needs
type SlackRequester interface {
	SearchMessagesContext(context.Context, string, slack.SearchParameters) (*slack.SearchMessages, error)
}

// SlackClient is api client for slack request
type SlackClient struct {
	Client SlackRequester
}

// NewSlackClient returns SlackClient object.
// arg httpclient for using custom http.Client. So this can be nil.
//
// For example,
// 	NewSlackClient('YOUR_TOKEN', urlfetch.Client(ctx))
func NewSlackClient(token string, httpclient *http.Client) *SlackClient {
	var opts []slack.Option

	if httpclient != nil {
		opts = []slack.Option{slack.OptionHTTPClient(httpclient)}
	}

	return &SlackClient{
		Client: slack.New(token, opts...),
	}
}

// CountQuery returns count of matches with query
func (sc *SlackClient) CountQuery(ctx context.Context, query string) (int, error) {
	res, err := sc.Client.SearchMessagesContext(ctx, query, slack.SearchParameters{})
	if err != nil {
		return 0, err
	}

	return res.Total, nil
}
