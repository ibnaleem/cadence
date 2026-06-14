package cmd

import (
	"fmt"
	"strconv"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <id|name>",
	Short: "Log a habit completion for today",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := util.InitDB()
		if err != nil {
			return err
		} // if
		defer db.Close()

		if err := util.SetupSchema(db); err != nil {
			return err
		} // if

		var habitID int
		var habitName string

		if id, err := strconv.Atoi(args[0]); err == nil {
			habitName, err = util.HabitNameByID(db, id)
			if err != nil {
				return fmt.Errorf("no habit with id %d — run `cadence list` to see habit IDs", id)
			} // if
			habitID = id
		} else {
			embedding, embErr := util.Embed(args[0])
			if embErr != nil {
				return fmt.Errorf("VOYAGE_API_KEY required for name-based lookup; use `cadence done <id>`")
			} // if

			habit, sim, err := util.FindSimilarHabit(db, embedding, util.SimilarityThresholdMatch)
			if err != nil {
				return err
			} // if
			if habit == nil {
				return fmt.Errorf("no habit matched \"%s\" (best similarity %.0f%%)", args[0], sim*100)
			} // if

			fmt.Printf("%s matched: %s\n", theme.Cyan("~"), theme.Bold(habit.Name))
			habitID = habit.ID
			habitName = habit.Name
		} // else

		if err := util.LogHabit(db, habitID); err != nil {
			return err
		} // if

		fmt.Println(theme.Green("✓") + " logged: " + theme.Bold(habitName))

		return nil
	}, // RunE
} // &cobra.Command
