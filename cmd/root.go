package cmd


import (
	"github.com/spf13/cobra"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/ibnaleem/cadence/internal/tui"
)

var rootCmd = &cobra.Command{
	Use: "cadence",
	Short: "👾 Your terminal-based habit tracker",

	RunE: func(cmd *cobra.Command, args []string) error {
    return tui.Render(cmd.UsageString())
	}, // RunE
} // &cobra.Command

func Execute() {

	err := rootCmd.Execute()

	util.CheckError(err)	

} // Execute