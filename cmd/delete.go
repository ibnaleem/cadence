package cmd

import (
	"fmt"
	"strconv"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Remove a habit and its history",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be a number — run `cadence list` to see habit IDs")
		} // if

		db, err := util.InitDB()
		if err != nil {
			return err
		} // if
		defer db.Close()

		if err := util.SetupSchema(db); err != nil {
			return err
		} // if

		name, err := util.HabitNameByID(db, id)
		if err != nil {
			return fmt.Errorf("no habit with id %d", id)
		} // if

		if err := util.DeleteHabit(db, id); err != nil {
			return err
		} // if

		fmt.Println(theme.Red("✗") + " deleted: " + name)

		return nil
	}, // RunE
} // &cobra.Command
