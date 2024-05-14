package cmd

import (
	"fmt"
	"github.com/pouspo/application-gateway/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	PolicyCommand = &cobra.Command{
		Use:   "insert-policy",
		Short: "Insert Dummy Policy",
		Run: func(cmd *cobra.Command, args []string) {
			err := internal.App.InsertDummyPolicy()
			if err != nil {
				logrus.Errorf("Could not insert dummy policy, error: %v", err)
				return
			}
		},
	}

	PolicyReadCommand = &cobra.Command{
		Use: "read-policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			var userName string
			if len(args) > 0 {
				userName = args[0]
			} else {
				fmt.Println("You have full access to your data")
				return nil
			}

			if userName == "" {
				userName = UserName
			}
			err := internal.App.ReadPolicy(userName)
			if err != nil {
				logrus.Errorf("Could not insert dummy policy, error: %v", err)
				return err
			}
			return nil
		},
	}
)
