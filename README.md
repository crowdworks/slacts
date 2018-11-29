# slacts

[![GoDoc](https://godoc.org/github.com/crowdworks/slacts?status.svg)](https://godoc.org/github.com/crowdworks/slacts)
[![codecov](https://codecov.io/gh/crowdworks/slacts/branch/master/graph/badge.svg)](https://codecov.io/gh/crowdworks/slacts)

A CLI tool for Slack statistics

## Getting Started

### Prerequisites

#### Generate Slack Token

To get Slack token, access here:
https://api.slack.com/custom-integrations/legacy-tokens

#### Set Slack token

Install [direnv](https://direnv.net/).

```bash
$ brew install direnv # MacOS
```

Copy .envrc from `.envrc.sample` and set Slack token.

```bash
$ cp .envrc.sample .envrc
$ vi .envrc

export SLACK_API_TOKEN=xxxxxxxxxxxxxxxxxxx
```

### Run

```bash
$ go run cmd/slacts/main.go slack count "in#general @channel"
result: 12
```

or

```bash
$ make install
go install github.com/crowdworks/slacts/cmd/slacts

$ slacts slack count "in#general @channel"
result: 12
```

## Synopsis

### `slacts slack count <slack_search_term>`

Returns the number of search result.
Please refer [Guide to search in Slack
](https://get.slack.help/hc/en-us/articles/202528808-Guide-to-search-in-Slack-).

## License

MIT
