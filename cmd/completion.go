package cmd

import (
	"fmt"
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
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = rootCmd.GenFishCompletion(os.Stdout, true)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating completion: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
