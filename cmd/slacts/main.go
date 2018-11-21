package main

import (
	"context"
	"fmt"
	"log"

	"github.com/crowdworks/slacts"
	"github.com/crowdworks/slacts/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slacts",
	Short: "a CLI tool for Slack statistics",
}

var slackCmd = &cobra.Command{
	Use:   "slack",
	Short: "commands for Slack",
}

// SlackCredentialSetter is interface for Slack credentials config
type SlackCredentialSetter interface {
	Token() string
}

var slackCountCmd = &cobra.Command{
	Use:     "count [query]",
	Short:   "count messages what matches query",
	Args:    cobra.ExactArgs(1),
	Example: "slacts slack count \"in:#general message\"",
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		var config SlackCredentialSetter = config.NewSlackClientConfig()

		if config.Token() == "" {
			log.Fatalln("unset Slack API Token")
		}

		client := slacts.NewSlackClient(config.Token(), nil)
		count, err := client.CountQuery(context.Background(), query)

		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("result: %d\n", count)
	},
}

func main() {
	slackCmd.AddCommand(slackCountCmd)

	rootCmd.AddCommand(slackCmd)
	rootCmd.Execute()
}
