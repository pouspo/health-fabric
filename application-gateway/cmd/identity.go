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
			var dob string
			if len(args) > 0 {
				dob = args[0]
			}
			if err := internal.App.RegisterAsPatient(UserName, dob); err != nil {
				fmt.Printf("Could not register as a patient, error: %v", err)
			}
		},
	}

	DiagnosisCmd = &cobra.Command{
		Use:   "create-diagnosis",
		Short: "Insert Diagnosis Data",
		Run: func(cmd *cobra.Command, args []string) {
			var userName string
			if len(args) > 0 {
				userName = args[0]
			}

			if userName == "" {
				userName = UserName
			}

			if err := internal.App.CreateDiagnosis(userName); err != nil {
				fmt.Printf("Could not create diagnosis, error: %v", err)
			}
		},
	}

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
