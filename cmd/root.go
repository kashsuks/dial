package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"dial/internal/store"
	"dial/internal/tracker"
)

var (
	db  *sql.DB
	trk *tracker.Tracker
)

var rootCmd = &cobra.Command{
	Use:   "dial",
	Short: "Dial - quick time tracking from the CLI and GUI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		path, err := store.DefaultPath()
		if err != nil {
			return err
		}
		db, err = store.Open(path)
		if err != nil {
			return err
		}
		trk = tracker.New(db)
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if db != nil {
			db.Close()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
