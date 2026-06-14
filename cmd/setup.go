package cmd

import (
	"fmt"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialise the local database",

	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := util.InitDB()
		if err != nil {
			return err
		} // if
		defer db.Close()

		if err := util.SetupSchema(db); err != nil {
			return err
		} // if

		path, _ := util.DBPath()
		fmt.Println(theme.Green("✓") + " database ready at " + theme.Gray(path))

		return nil
	}, // RunE
} // &cobra.Command
