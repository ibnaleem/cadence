package cmd

import (
	"fmt"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all habits",

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
			fmt.Println(theme.Gray("No habits yet. Run `cadence add <name>` to get started."))
			return nil
		} // if

		for _, h := range habits {
			freq := theme.Gray("[" + h.Frequency + "]")
			desc := ""
			if h.Description != "" {
				desc = "  " + theme.Gray(h.Description)
			} // if
			fmt.Printf("%s %s%s\n", theme.Cyan(fmt.Sprintf("%d.", h.ID)), theme.Bold(h.Name)+" "+freq, desc)
		} // for

		return nil
	}, // RunE
} // &cobra.Command
