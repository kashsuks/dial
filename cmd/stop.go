package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"dial/internal/tracker"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the currently running session",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := trk.Stop()
		if err != nil {
			if errors.Is(err, tracker.ErrNoRunningSession) {
				fmt.Println("No session is currently running.")
				return nil
			}
			return err
		}
		dur := s.EndedAt.Sub(s.StartedAt)
		fmt.Printf("Stopped: %s (%s)\n", s.Task, dur.Round(1e9))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
