package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
	appconfig "github.com/planitaicojp/freee-cli/internal/config"
	"github.com/planitaicojp/freee-cli/internal/output"
)

// Cmd is the config command group.
var Cmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

func init() {
	Cmd.AddCommand(showCmd)
	Cmd.AddCommand(setCmd)
	Cmd.AddCommand(pathCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := appconfig.Load()
		if err != nil {
			return err
		}
		return output.New(cmdutil.GetFormat(cmd)).Format(os.Stdout, cfg)
	},
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cmdutil.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := appconfig.Load()
		if err != nil {
			return err
		}

		key, value := args[0], args[1]
		switch key {
		case "format":
			cfg.Defaults.Format = value
		default:
			return fmt.Errorf("unknown config key: %s (available: format)", key)
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Set %s = %s\n", key, value)
		return nil
	},
}

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print config directory path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(appconfig.DefaultConfigDir())
	},
}
