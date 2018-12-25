package cmd

import (
	"github.com/crowdworks/slacts/config"
	"github.com/spf13/cobra"
)

var (
	slackConfig *config.SlackClientConfig
)

// RootCmd is entry point of commands
// sub command defines at init() function
var RootCmd = &cobra.Command{
	Use:   "slacts",
	Short: "a CLI tool for Slack statistics",
}
