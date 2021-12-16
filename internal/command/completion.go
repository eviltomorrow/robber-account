package command

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(robber-account completion bash)

# To load completions for each session, execute once:
Linux:
  $ robber-account completion bash > /etc/bash_completion.d/robber-account
MacOS:
  $ robber-account completion bash > /usr/local/etc/bash_completion.d/robber-account

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ robber-account completion zsh > "${fpath[1]}/_robber-account"

# You will need to start a new shell for this setup to take effect.

Fish:

$ robber-account completion fish | source

# To load completions for each session, execute once:
$ robber-account completion fish > ~/.config/fish/completions/robber-account.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
