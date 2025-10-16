package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(desk completion bash)

  To load completions for each session, execute once:
  # Linux:
  $ desk completion bash > /etc/bash_completion.d/desk
  # macOS:
  $ desk completion bash > /usr/local/etc/bash_completion.d/desk

Zsh:

  $ source <(desk completion zsh)

  To load completions for each session, execute once:
  $ desk completion zsh > "${fpath[1]}/_desk"

Fish:

  $ desk completion fish | source

  To load completions for each session, execute once:
  $ desk completion fish > ~/.config/fish/completions/desk.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			rootCmd.GenFishCompletion(os.Stdout, true)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
