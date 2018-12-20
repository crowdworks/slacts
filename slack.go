package slacts

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

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
func (sc *SlackClient) CountQuery(ctx context.Context, query *SlackQuery) (int, error) {
	res, err := sc.Client.SearchMessagesContext(ctx, query.String(), slack.SearchParameters{})
	if err != nil {
		return 0, err
	}

	return res.Total, nil
}

// SlackQuery for slack message search
type SlackQuery string

// TODO: enable to parse other format: for example 2018-01-01
var queryDateRegexp = regexp.MustCompile(`on:(\d{4}\/\d{1,2}\/\d{1,2})`)

// NewSlackQuery is initializer of Slack query
func NewSlackQuery(query string) *SlackQuery {
	sq := SlackQuery(query)
	return &sq
}

// Date returns date
func (q *SlackQuery) Date() (*time.Time, error) {
	matches := queryDateRegexp.FindAllStringSubmatch(string(*q), -1)

	if len(matches) == 0 || len(matches[0]) < 2 {
		return nil, errors.New("did not find date")
	}

	date, err := time.Parse("2006/01/02", matches[0][1])
	if err != nil {
		return nil, err
	}

	return &date, nil
}

// String .
func (q *SlackQuery) String() string {
	return string(*q)
}
