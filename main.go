package main

import (
	"fmt"
	"github.com/ibnaleem/cadence/cmd"
	"github.com/ibnaleem/cadence/internal/util"
	"github.com/ibnaleem/cadence/internal/theme"
)

func main() {

	fmt.Println()
	fmt.Println(theme.Bold(theme.Cyan("cadence")) + theme.Gray(" · ") + theme.Gray(util.GetRandTitlePhrase()))
	fmt.Println()

	cmd.Execute()
} // main