package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	logProject string
	logTags string
	logDuration string
	logAt string
)

var logCmd = &cobra.Command{
	Use: "log [task]",
	Short: "Log a completed session retroactively",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dur, err := time.ParseDuration(logDuration)
		if err != nil {
			return fmt.Errorf("invalid --duration: %w", err)
		}

		end := time.Now()
		if logAt != "" {
			end, err = time.Parse("15:04", logAt)
			if err != nil {
				return fmt.Errorf("invalid --at (use HH:MM): %w", err)
			}
			now := time.Now()
			end = time.Date(now.Year(), now.Month(), now.Day(), end.Hour(), end.Minute(), 0, 0, now.Location())
		}
		start := end.Add(-dur)

		s, err := trk.Log(args[0], logProject, logTags, "cli", start, end)
		if err != nil {
			return err
		}
		fmt.Printf("Logged: %s (%s, %s - %s)\n", s.Task, dur, start.Format("15:04"), end.Format("15:04"))
		return nil
	},
}

func init() {
	logCmd.Flags().StringVar(&logProject, "project", "", "projet name")
	logCmd.Flags().StringVar(&logTags, "tag", "", "comma-seperated tags")
	logCmd.Flags().StringVar(&logDuration, "duration", "", "duration, e.g. 30m, 1h15m (required)")
	logCmd.Flags().StringVar(&logAt, "at", "", "end time HH:MM (defaults to now)")
	logCmd.MarkFlagRequired("duration")
	rootCmd.AddCommand(logCmd)
}
