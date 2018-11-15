package slacts_test

import (
	"context"
	"testing"

	"github.com/crowdworks/slacts"
	"github.com/nlopes/slack"
)

type testSlackClient struct{}

func (tsc *testSlackClient) SearchMessagesContext(ctx context.Context, q string, params slack.SearchParameters) (*slack.SearchMessages, error) {
	sm := &slack.SearchMessages{
		Total: 27,
	}

	return sm, nil
}

func TestSlackClient_CountQuery(t *testing.T) {
	ctx := context.Background()

	sc := slacts.SlackClient{
		Client: new(testSlackClient),
	}

	count, err := sc.CountQuery(ctx, "in:general channel")
	if err != nil {
		t.Error(err)
	}

	if count != 27 {
		t.Errorf("query count expected 27 but actual %d", count)
	}
}
