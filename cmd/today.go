package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use: "today",
	Short: "Show today's tracked sessions and total time",
	RunE: func(cmd *cobra.Command, args []string) error {
		now := time.Now()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		rows, err := db.Query(
			`SELECT task, project, started_at, ended_at FROM sessions
			 WHERE started_at >= ? ORDER BY started_at ASC`,
			startOfDay,
		)
		if err != nil {
			return err
		}
		defer rows.Close()

		var total time.Duration
		count := 0
		for rows.Next() {
			var task, project string
			var started time.Time
			var ended *time.Time
			if err := rows.Scan(&task, &project, &started, &ended); err != nil {
				return err
			}
			end := now
			label := "(running)"
			if ended != nil {
				end = *ended
				label = end.Format("15:04")
			}
			dur := end.Sub(started)
			total += dur
			count++
			fmt.Printf("%s - %s %-20s %s\n", started.Format("15:04"), label, task, dur.Round(1e9))
		}

		if count == 0 {
			fmt.Println("No sesions logged today.")
			return nil
		}
		fmt.Printf("\nTotal: %s\n", total.Round(1e9))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(todayCmd)
}
