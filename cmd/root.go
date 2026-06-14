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

		weekly, err := util.WeeklyLogs(db)
		if err != nil {
			return err
		} // if

		totalWeek := 0
		for _, c := range weekly {
			totalWeek += c
		} // for

		fmt.Println()
		fmt.Println("  " + theme.Bold("This week") + "  " + theme.Gray(weekLabel(totalWeek)))
		fmt.Println()

		for _, h := range habits {
			count := weekly[h.ID]
			dots := weekDots(count, 7)
			fmt.Printf("  %s  %-20s  %s  %s\n",
				dots,
				h.Name,
				theme.Cyan(fmt.Sprintf("%dx", count)),
				theme.Gray("["+h.Frequency+"]"),
			)
		} // for

		fmt.Println()
		return nil
	}, // RunE
} // &cobra.Command

func weekDots(count, max int) string {
	if count > max {
		count = max
	} // if
	return theme.Green(strings.Repeat("●", count)) + theme.Gray(strings.Repeat("●", max-count))
} // weekDots

func weekLabel(total int) string {
	switch {
	case total == 0:
		return "nothing logged yet — let's go!"
	case total < 5:
		return "warming up 🔥"
	case total < 10:
		return "building momentum ⚡"
	case total < 20:
		return "on a roll 🚀"
	default:
		return "unstoppable 🏆"
	} // switch
} // weekLabel

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
