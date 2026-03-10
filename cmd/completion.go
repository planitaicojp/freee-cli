package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/freee-cli/cmd/cmdutil"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate shell completion script for the specified shell.

Examples:
  # Bash
  freee completion bash > /etc/bash_completion.d/freee

  # Zsh
  freee completion zsh > "${fpath[1]}/_freee"

  # Fish
  freee completion fish > ~/.config/fish/completions/freee.fish`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cmdutil.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return nil
		}
	},
}
