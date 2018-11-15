package slacts_test

import (
	"context"
	"errors"
	"testing"

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

			count, err := sc.CountQuery(ctx, "in:general channel")

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
