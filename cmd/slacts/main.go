package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/crowdworks/slacts"
	"github.com/crowdworks/slacts/config"
	"github.com/crowdworks/slacts/report"
	"github.com/spf13/cobra"
)

var (
	slackConfig *config.SlackClientConfig
)

func init() {
	slackConfig = config.NewSlackClientConfig()

	rootCmd.AddCommand(newTaskCmd(), newSlackCmd())
}

var rootCmd = &cobra.Command{
	Use:   "slacts",
	Short: "a CLI tool for Slack statistics",
}

func newTaskCmd() *cobra.Command {
	type option struct {
		file string
	}

	var opt option
	cmd := &cobra.Command{
		Use:     "task",
		Short:   "exec tasks",
		Example: "slacts task slacts.yml",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if opt.file == "" {
				return errors.New("task definition file is required")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := config.ReadYaml(opt.file, taskNameFilters(args)...)
			if err != nil {
				log.Fatal(err)
			}

			for _, task := range *tasks {
				switch task.Kind {
				case "count":
					query := slacts.NewSlackQuery(task.Query)

					if slackConfig.Token() == "" {
						log.Fatalln("unset Slack API Token")
					}

					client := slacts.NewSlackClient(slackConfig.Token(), nil)
					count, err := client.CountQuery(context.Background(), query)

					if err != nil {
						log.Fatalln(err)
					}

					log.Printf("name: %s, kind: %s, result: %d\n", task.Name, task.Kind, count)

					if task.DoesSendDatadog() {
						ddConfig := config.NewDatadogConfig()

						if ddConfig.AppKey() == "" || ddConfig.APIKey() == "" {
							log.Fatalln("unset Datadog API credentials. Need both API key and application key")
						}

						ddClient := report.NewDatadogClient(ddConfig.APIKey(), ddConfig.AppKey(), nil)

						// TODO: add lacked properties. for example, Unit, Host and etc...
						metrics := []report.DatadogMetric{
							{
								Metric: &task.Datadog.Metric,
								Points: []report.DatadogDataPoint{
									{
										float64Pointer(float64(time.Now().Unix())),
										float64Pointer(float64(count)),
									},
								},
								Tags: task.Datadog.Tags,
							},
						}

						if err := ddClient.PostMetrics(metrics); err != nil {
							log.Fatal(err)
						}
						log.Println("send metric to datadog successfully")
					}
				}
			}
		},
	}

	cmd.PersistentFlags().StringVarP(&opt.file, "file", "f", "", "task definition file path")

	return cmd
}

func taskNameFilters(names []string) []config.ReadYamlOption {
	var opts []config.ReadYamlOption
	for _, name := range names {
		opts = append(opts, config.OptionNameFilter(name))
	}
	return opts
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

func main() {
	_ = rootCmd.Execute()
}
