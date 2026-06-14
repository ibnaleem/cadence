package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/ibnaleem/cadence/internal/theme"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cadence",
	Short: "👾 Your terminal-based habit tracker",

	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := util.InitDB()
		if err != nil {
			return err
		} // if
		defer db.Close()

		if err := util.SetupSchema(db); err != nil {
			return err
		} // if

		habits, err := util.TodayStatus(db)
		if err != nil {
			return err
		} // if

		if len(habits) == 0 {
			fmt.Println("  " + theme.Gray("No habits yet — run `cadence add <name>` to get started."))
			fmt.Println()
			return nil
		} // if

		done := 0
		for _, h := range habits {
			if h.DoneToday {
				done++
			} // if
		} // for

		pct := (done * 100) / len(habits)

		fmt.Println("  " + theme.Gray(time.Now().Format("Monday, January 2")))
		fmt.Println()
		fmt.Printf("  %s  %s\n",
			progressBar(done, len(habits), 28),
			theme.Bold(fmt.Sprintf("%d/%d", done, len(habits)))+theme.Gray(fmt.Sprintf(" · %d%%", pct)),
		)
		fmt.Println()

		for _, h := range habits {
			if h.DoneToday {
				fmt.Printf("  %s  %s  %s\n", theme.Green("✓"), theme.Green(theme.Bold(h.Name)), theme.Gray("["+h.Frequency+"]"))
			} else {
				fmt.Printf("  %s  %s  %s\n", theme.Gray("○"), theme.Gray(h.Name), theme.Gray("["+h.Frequency+"]"))
			} // if
		} // for

		fmt.Println()
		return nil
	}, // RunE
} // &cobra.Command

func progressBar(done, total, width int) string {
	filled := (done * width) / total
	return theme.Green(strings.Repeat("█", filled)) + theme.Gray(strings.Repeat("░", width-filled))
} // progressBar

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(doneCmd)
} // init

func Execute() {
	err := rootCmd.Execute()
	util.CheckError(err)
} // Execute
