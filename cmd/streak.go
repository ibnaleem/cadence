package cmd

import (
	"fmt"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var streakCmd = &cobra.Command{
	Use:   "streak",
	Short: "Show current streak for each habit",

	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := util.InitDB()
		if err != nil {
			return err
		} // if
		defer db.Close()

		if err := util.SetupSchema(db); err != nil {
			return err
		} // if

		habits, err := util.ListHabits(db)
		if err != nil {
			return err
		} // if
		if len(habits) == 0 {
			fmt.Println(theme.Gray("  No habits yet."))
			return nil
		} // if

		streaks, err := util.AllStreaks(db)
		if err != nil {
			return err
		} // if

		fmt.Println()
		for _, h := range habits {
			s := streaks[h.ID]
			var indicator string
			if s > 0 {
				indicator = fmt.Sprintf("🔥 %s", theme.Bold(fmt.Sprintf("%d-day streak", s)))
			} else {
				indicator = theme.Gray("no streak")
			} // if
			fmt.Printf("  %-24s %s\n", h.Name, indicator)
		} // for
		fmt.Println()

		return nil
	}, // RunE
} // &cobra.Command
