package plugin

import (
	"errors"
	"fmt"
	"github.com/kaytu-io/kaytu/pkg/plugin"
	"github.com/kaytu-io/kaytu/pkg/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use: "install",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := plugin.New()
		err := manager.StartServer()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return errors.New("please provide plugin path")
		}

		token := utils.ReadStringFlag(cmd, "token")
		err = manager.Install(args[0], token)
		if err != nil {
			fmt.Printf("failed to install plugin due to %v\n", err)
			return err
		}
		return nil
	},
}
