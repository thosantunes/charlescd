package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func newCmdInstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Charles by .yaml configuration file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a .yaml configuation file")
			}
			return nil
		},
		Run: cmdInstallRun,
	}

	return cmd
}
