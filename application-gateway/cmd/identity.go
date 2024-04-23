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

	DiagnosisCmd = &cobra.Command{
		Use:   "diagnosis",
		Short: "Manage Diagnosis",
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

func init() {
	DiagnosisCmd.AddCommand(&cobra.Command{
		Use:   "insert",
		Short: "Insert Diagnosis Data",
		Long:  "Insert diagnosis data pass them as pair of key value like [key1 value1 key2 value2 ... keyn valuen]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := internal.App.InsertDiagnosisData(args...); err != nil {
				fmt.Printf("Could not create diagnosis, error: %v", err)
				return err
			}
			return nil
		},
	})
}
