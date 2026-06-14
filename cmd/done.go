package cmd

import (
	"fmt"
	"strconv"
	"time"

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
				return fmt.Errorf(
					"ollama not reachable — name-based lookup requires it.\n\n" +
					"  Install:  https://ollama.com/download\n" +
					"  Pull model: ollama pull embeddinggemma\n\n" +
					"Or use a habit ID instead: cadence done <id>",
				)
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

		dateStr, _ := cmd.Flags().GetString("date")
		if dateStr != "" {
			if _, err := time.Parse("2006-01-02", dateStr); err != nil {
				return fmt.Errorf("invalid date %q — use YYYY-MM-DD format", dateStr)
			} // if
		} // if

		inserted, err := util.LogHabit(db, habitID, dateStr)
		if err != nil {
			return err
		} // if

		if dateStr != "" {
			if inserted {
				fmt.Println(theme.Green("✓") + " logged " + theme.Bold(habitName) + theme.Gray(" for "+dateStr))
			} else {
				fmt.Println(theme.Yellow("!") + " already logged " + theme.Bold(habitName) + theme.Gray(" for "+dateStr))
			} // if
		} else {
			if inserted {
				fmt.Println(theme.Green("✓") + " logged: " + theme.Bold(habitName))
			} else {
				fmt.Println(theme.Yellow("!") + " already logged today: " + theme.Bold(habitName))
			} // if
		} // if

		return nil
	}, // RunE
} // &cobra.Command

func init() {
	doneCmd.Flags().StringP("date", "d", "", "backfill date (YYYY-MM-DD)")
} // init
