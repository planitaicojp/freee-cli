package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("freee-cli version %s by PlanitAI Co., Ltd.\n", version)
		fmt.Println("This is an unofficial tool and is not affiliated with or endorsed by freee K.K.")
	},
}
