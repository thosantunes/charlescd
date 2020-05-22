package cli

import (
	"github.com/spf13/cobra"
)

type UseCases interface {
	Execute()
}

type CLI struct {
	rootCmd *cobra.Command
}

func NewCLI() UseCases {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(newCmdInstall())
	return &CLI{rootCmd}
}

func (cli *CLI) Execute() {
	cli.rootCmd.Execute()
}
