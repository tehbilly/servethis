package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "What version of 'servethis' are you using?",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("servethis version 0.0.1")
	},
}
