package cli // import "github.com/automaticserver/lxe/cli"

import (
	"github.com/spf13/cobra"
)

func initCmpl() {
	rootCmd.AddCommand(cmplCmd)

	cmplCmd.AddCommand(cmplBashCmd)
	cmplCmd.AddCommand(cmplZshCmd)
	cmplCmd.AddCommand(cmplPwrshCmd)
}

var cmplCmd = &cobra.Command{
	Use:          "completion",
	Aliases:      []string{"compl"},
	Short:        "Generate a completion script",
	Args:         cobra.NoArgs,
	Run:          func(*cobra.Command, []string) {},
	SilenceUsage: true,
}

var cmplBashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generates bash completion script",
	Long: `To load completion run

. <(... generate completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(... generate completion bash)
`,
	RunE: cmplBashCmdRunE,
}

func cmplBashCmdRunE(cmd *cobra.Command, args []string) error {
	return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
}

var cmplZshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generates zsh completion script",
	RunE:  cmplZshCmdRunE,
}

func cmplZshCmdRunE(cmd *cobra.Command, args []string) error {
	return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
}

var cmplPwrshCmd = &cobra.Command{
	Use:   "powershell",
	Short: "Generates powershell completion script",
	RunE:  cmplPwrshCmdRunE,
}

func cmplPwrshCmdRunE(cmd *cobra.Command, args []string) error {
	return cmd.Root().GenPowerShellCompletion(cmd.OutOrStdout())
}
