package slacts_test

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/crowdworks/slacts"
	"github.com/nlopes/slack"
)

type testSlackClient struct {
	hasError bool
}

func (tsc *testSlackClient) SearchMessagesContext(ctx context.Context, q string, params slack.SearchParameters) (*slack.SearchMessages, error) {

	if tsc.hasError {
		return nil, errors.New("something error occurred")
	}

	sm := &slack.SearchMessages{
		Total: 27,
	}

	return sm, nil
}

func TestSlackClient_CountQuery(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		client               *testSlackClient
		expectError          bool
		expectedErrorMessage string
	}{
		"normal": {
			client:      &testSlackClient{},
			expectError: false,
		},
		"has something error": {
			client:               &testSlackClient{hasError: true},
			expectError:          true,
			expectedErrorMessage: "something error occurred",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			sc := slacts.SlackClient{
				Client: c.client,
			}

			count, err := sc.CountQuery(ctx, slacts.NewSlackQuery("in:#general channel"))

			if c.expectError {
				if err == nil {
					t.Error("expected to occur error but no errors occurred")
				}

				if c.expectedErrorMessage != err.Error() {
					t.Errorf("unexpected error messages: expected '%s' actual '%s'", c.expectedErrorMessage, err)
				}

				return
			}

			if err != nil {
				t.Error(err)
			}

			if count != 27 {
				t.Errorf("query count expected 27 but actual %d", count)
			}
		})
	}
}

func TestNewSlackClient(t *testing.T) {
	type args struct {
		token      string
		httpclient *http.Client
	}

	cases := map[string]struct {
		args args
		want *slacts.SlackClient
	}{
		"default": {
			args: args{
				token:      "aaa",
				httpclient: nil,
			},
			want: &slacts.SlackClient{
				Client: slack.New("aaa"),
			},
		},
		"custom client": {
			args: args{
				token: "bbb",
				httpclient: &http.Client{
					Transport:     nil,
					CheckRedirect: nil,
					Jar:           nil,
					Timeout:       30 * time.Second,
				},
			},
			want: &slacts.SlackClient{
				Client: slack.New("bbb", slack.OptionHTTPClient(&http.Client{
					Transport:     nil,
					CheckRedirect: nil,
					Jar:           nil,
					Timeout:       30 * time.Second,
				})),
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := slacts.NewSlackClient(c.args.token, c.args.httpclient); !reflect.DeepEqual(got, c.want) {
				t.Errorf("NewSlackClient() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestSlackQuery_Date(t *testing.T) {
	cases := map[string]struct {
		q       *slacts.SlackQuery
		want    time.Time
		wantErr bool
	}{
		"only date": {
			q:       slacts.NewSlackQuery("on:2018/02/28"),
			want:    time.Date(2018, 2, 28, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		"with other query": {
			q:       slacts.NewSlackQuery("in:#general on:2018/02/28 @channel"),
			want:    time.Date(2018, 2, 28, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		"wrong date format": {
			q:       slacts.NewSlackQuery("in:#general on:2018/22/28 @channel"),
			wantErr: true,
		},
		"no date in query": {
			q:       slacts.NewSlackQuery("in:#general @channel"),
			wantErr: true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := c.q.Date()
			if (err != nil) != c.wantErr {
				t.Errorf("SlackQuery.Date() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if c.wantErr {
				return
			}

			if !reflect.DeepEqual(got, &c.want) {
				t.Errorf("SlackQuery.Date() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestSlackQuery_String(t *testing.T) {
	cases := map[string]struct {
		q    *slacts.SlackQuery
		want string
	}{
		"normal": {
			q:    slacts.NewSlackQuery("in:#general on:2018/02/28 @channel"),
			want: "in:#general on:2018/02/28 @channel",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := c.q.String(); got != c.want {
				t.Errorf("SlackQuery.String() = %v, want %v", got, c.want)
			}
		})
	}
}
