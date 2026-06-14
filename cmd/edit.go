package cmd

import (
	"fmt"
	"strconv"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <id>",
	Short: "Update a habit's name or description",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("id must be a number — run `cadence list` to see habit IDs")
		} // if

		if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("description") {
			return fmt.Errorf("provide at least --name or --description")
		} // if

		db, err := util.InitDB()
		if err != nil {
			return err
		} // if
		defer db.Close()

		if err := util.SetupSchema(db); err != nil {
			return err
		} // if

		habit, err := util.GetHabit(db, id)
		if err != nil {
			return fmt.Errorf("no habit with id %d", id)
		} // if

		var newEmbedding []float32
		if cmd.Flags().Changed("name") {
			habit.Name, _ = cmd.Flags().GetString("name")
			emb, err := util.Embed(habit.Name)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not update embedding (%v); name-based lookup may be less accurate\n", err)
			} else {
				newEmbedding = emb
			}
		} // if
		if cmd.Flags().Changed("description") {
			habit.Description, _ = cmd.Flags().GetString("description")
		} // if

		if err := util.UpdateHabit(db, id, habit.Name, habit.Description, newEmbedding); err != nil {
			return err
		} // if

		fmt.Println(theme.Green("✓") + " updated: " + theme.Bold(habit.Name))

		return nil
	}, // RunE
} // &cobra.Command

func init() {
	editCmd.Flags().StringP("name", "n", "", "new name")
	editCmd.Flags().StringP("description", "d", "", "new description")
} // init
