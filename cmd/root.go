package cmd


import (
	"github.com/spf13/cobra"
	"github.com/ibnaleem/cadence/internal/util"
)

var rootCmd = &cobra.Command{
	Use: "cadence",
	Short: "👾 Your terminal-based habit tracker",

	RunE: func(cmd *cobra.Command, args []string) error {
    return cmd.Help()
	}, // RunE
} // &cobra.Command

func init() {
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
} // init

func Execute() {

	err := rootCmd.Execute()

	util.CheckError(err)

} // Execute