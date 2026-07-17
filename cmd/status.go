package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"dial/internal/tracker"
)

var statusCmd = &cobra.Command{
	Use: "status",
	Short: "Show the currently running session, if any",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := trk.Current()
		if err != nil {
			if errors.Is(err, tracker.ErrNoRunningSession) {
				fmt.Println("Nothing is currently being tracked.")
				return nil
			}
			return err
		}
		elapsed := time.Since(s.StartedAt).Round(1e9)
		fmt.Printf("Running: %s (%s elapsed)\n", s.Task, elapsed)
		if s.Project != "" {
			fmt.Printf(" project: %s\n", s.Project)
		}
		if s.Tags != "" {
			fmt.Printf(" tags: %s\n", s.Tags)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
