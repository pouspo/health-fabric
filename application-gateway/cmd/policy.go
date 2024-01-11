package cmd

import (
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
)
