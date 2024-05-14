package cmd

import (
	"fmt"
	"github.com/pouspo/application-gateway/internal"
	"github.com/spf13/cobra"
)

var (
	ReadCmd = &cobra.Command{
		Use:   "read",
		Short: "Read user data",
		Run: func(cmd *cobra.Command, args []string) {
			var userName string
			if len(args) > 0 {
				userName = args[0]
			}

			if userName == "" {
				userName = UserName
			}

			if err := internal.App.ReadUserData(userName); err != nil {
				fmt.Printf("Could not create diagnosis, error: %v", err)
			}
		},
	}
)
