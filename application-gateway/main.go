package main

import (
	"fmt"
	"github.com/pouspo/application-gateway/cmd"
	"github.com/pouspo/application-gateway/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gateway",
		Short: "A Cli tool to interact with ABAC health data blockchain",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(func() {
		app, err := internal.NewApplication(cmd.UserName)
		if err != nil {
			panic(fmt.Errorf("could not create applicatoin: error: %v", err))
		}

		internal.App = app
	})
	rootCmd.PersistentFlags().StringVarP(&cmd.UserName, "user", "u", "alpha", "User Name: alpha, beta or gama")

	rootCmd.AddCommand(cmd.IdentityCommand)
	rootCmd.AddCommand(cmd.PolicyCommand)
	rootCmd.AddCommand(cmd.RegisterCmd)
	rootCmd.AddCommand(cmd.ReadCmd)
	rootCmd.AddCommand(cmd.CSVDiagnosisCmd)
	rootCmd.AddCommand(cmd.ListenEventCmd)

}

func main() {
	if err := Execute(); err != nil {
		logrus.Errorf("Error: %v", err)
	}
}
