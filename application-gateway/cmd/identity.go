package cmd

import (
	"fmt"
	"github.com/pouspo/application-gateway/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	UserName        string
	IdentityCommand = &cobra.Command{
		Use:   "id",
		Short: "Get user id",
		Run: func(cmd *cobra.Command, args []string) {
			id, err := internal.App.GetUserId()
			if err != nil {
				logrus.Errorf("Could not get userid, error: %v", err)
				return
			}

			fmt.Println("UserId: ", id)
		},
	}

	RegisterCmd = &cobra.Command{
		Use:   "register",
		Short: "Register as a patient",
		Run: func(cmd *cobra.Command, args []string) {
			if err := internal.App.RegisterAsPatient(UserName, ""); err != nil {
				fmt.Printf("Could not register as a patient, error: %v", err)
			}
		},
	}

	ListenEventCmd = &cobra.Command{
		Use:   "listen",
		Short: "Listen Block Events",
		Run: func(cmd *cobra.Command, args []string) {

			if err := internal.App.ListenBlockEvents(); err != nil {
				fmt.Printf("Could not listen, error: %v", err)
			}
		},
	}
)
