package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo <id>",
	Short: "Remove today's log for a habit",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be an integer — run `cadence list` to see habit IDs")
		} // if

		db, err := util.InitDB()
		if err != nil {
			return err
		} // if
		defer db.Close()

		if err := util.SetupSchema(db); err != nil {
			return err
		} // if

		habitName, err := util.HabitNameByID(db, id)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no habit with id %d — run `cadence list` to see habit IDs", id)
		} else if err != nil {
			return fmt.Errorf("looking up habit %d: %w", id, err)
		} // if

		removed, err := util.UnlogHabit(db, id)
		if err != nil {
			return err
		} // if

		if removed {
			fmt.Println(theme.Yellow("↩") + " undone: " + theme.Bold(habitName))
		} else {
			fmt.Println(theme.Gray("nothing to undo for " + theme.Bold(habitName) + " today"))
		} // if

		return nil
	}, // RunE
} // &cobra.Command
