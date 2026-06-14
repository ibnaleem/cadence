package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new habit",
	Args:  cobra.ExactArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		description, _ := cmd.Flags().GetString("description")
		frequency, _ := cmd.Flags().GetString("frequency")

		db, err := util.InitDB()
		if err != nil {
			return err
		} // if
		defer db.Close()

		if err := util.SetupSchema(db); err != nil {
			return err
		} // if

		embedding, _ := util.Embed(args[0])

		if embedding != nil {
			similar, sim, err := util.FindSimilarHabit(db, embedding, util.SimilarityThresholdWarn)
			if err != nil {
				return err
			} // if
			if similar != nil {
				fmt.Printf(
					"%s \"%s\" is %.0f%% similar to existing habit \"%s\". Add anyway? [y/N] ",
					theme.Yellow("!"), args[0], sim*100, theme.Bold(similar.Name),
				)
				ans, _ := bufio.NewReader(os.Stdin).ReadString('\n')
				if !strings.EqualFold(strings.TrimSpace(ans), "y") {
					fmt.Println(theme.Gray("aborted."))
					return nil
				} // if
			} // if
		} // if

		if err := util.AddHabit(db, args[0], description, frequency, embedding); err != nil {
			return err
		} // if

		fmt.Println(theme.Green("✓") + " habit added: " + theme.Bold(args[0]))

		return nil
	}, // RunE
} // &cobra.Command

func init() {
	addCmd.Flags().StringP("description", "d", "", "optional description")
	addCmd.Flags().StringP("frequency", "f", "daily", "daily or weekly")
} // init
