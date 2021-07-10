// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
	"io"
)

// NewCompletionCmd completionCmd represents the completion command
func NewCompletionCmd(out io.Writer) *cobra.Command {

	var completionCmd = &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

  $ source <(blackbean completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ blackbean completion bash > /etc/bash_completion.d/blackbean
  # macOS:
  $ blackbean completion bash > /usr/local/etc/bash_completion.d/blackbean

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ blackbean completion zsh > "${fpath[1]}/_blackbean"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ blackbean completion fish | source

  # To load completions for each session, execute once:
  $ blackbean completion fish > ~/.config/fish/completions/blackbean.fish

PowerShell:

  PS> blackbean completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> blackbean completion powershell > blackbean.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(out)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(out)
			case "fish":
				_ = cmd.Root().GenFishCompletion(out, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletionWithDesc(out)
			}
		},
	}
	return completionCmd
}

// Function to disable file completion
func noCompletions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}
