package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/crowdworks/slacts/config"

	"github.com/crowdworks/slacts"
	"github.com/spf13/cobra"
)

func init() {
	slackConfig = config.NewSlackClientConfig()
	RootCmd.AddCommand(newSlackCmd())
}

// newSlackCmd
func newSlackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slack",
		Short: "commands for Slack",
	}

	cmd.AddCommand(newSlackCountCmd())

	return cmd
}

func newSlackCountCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "count [query]",
		Short:   "count messages what matches query",
		Args:    cobra.ExactArgs(1),
		Example: "slacts slack count \"in:#general message\"",
		Run: func(cmd *cobra.Command, args []string) {
			query := slacts.NewSlackQuery(args[0])

			if slackConfig.Token() == "" {
				log.Fatalln("unset Slack API Token")
			}

			client := slacts.NewSlackClient(slackConfig.Token(), nil)
			count, err := client.CountQuery(context.Background(), query)

			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("result: %d\n", count)
		},
	}
}

func float64Pointer(v float64) *float64 {
	return &v
}
