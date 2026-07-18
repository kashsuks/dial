package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	startProject string
	startTags    string
)

var startCmd = &cobra.Command{
	Use:   "start [task]",
	Short: "Start tracking a task (stops any currently running session)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := trk.Start(args[0], startProject, startTags, "cli")
		if err != nil {
			return err
		}
		fmt.Printf("Started: %s (started at %s)\n", s.Task, s.StartedAt.Format("15:04:05"))
		return nil
	},
}

func init() {
	startCmd.Flags().StringVar(&startProject, "project", "", "project name")
	startCmd.Flags().StringVar(&startTags, "tag", "", "comma-seperated tags")
	rootCmd.AddCommand(startCmd)
}
