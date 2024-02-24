package cmd

import (
	"fmt"
	"github.com/pouspo/application-gateway/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
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

	DummyDiagnosisCmd = &cobra.Command{
		Use:   "create-dummy-diagnosis",
		Short: "Insert Diagnosis Data",
		Run: func(cmd *cobra.Command, args []string) {
			var userName string
			if len(args) > 0 {
				userName = args[0]
			}

			if userName == "" {
				userName = UserName
			}

			if err := internal.App.CreateDummyDiagnosis(userName); err != nil {
				fmt.Printf("Could not create diagnosis, error: %v", err)
			}
		},
	}

	CSVDiagnosisCmd = &cobra.Command{
		Use:   "insert-csv-diagnosis",
		Short: "Insert Diagnosis Data From CSV",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 9 {
				fmt.Println("Not enough arguments")
				return
			}
			var (
				pregnancies, glucose, bloodPressure, skinThickness, insulin, BMI int64
				Age, Outcome                                                     int64
				DiabetesPedigreeFunction                                         float64
			)

			pregnancies, _ = strconv.ParseInt(args[0], 10, 32)
			glucose, _ = strconv.ParseInt(args[1], 10, 32)
			bloodPressure, _ = strconv.ParseInt(args[2], 10, 32)
			skinThickness, _ = strconv.ParseInt(args[3], 10, 32)
			insulin, _ = strconv.ParseInt(args[4], 10, 32)
			BMI, _ = strconv.ParseInt(args[5], 10, 32)
			DiabetesPedigreeFunction, _ = strconv.ParseFloat(args[6], 32)
			Age, _ = strconv.ParseInt(args[7], 10, 32)
			Outcome, _ = strconv.ParseInt(args[8], 10, 32)

			if err := internal.App.InsertDiagnosisFromPimaDiabetesDataset(
				pregnancies,
				glucose,
				bloodPressure,
				skinThickness,
				insulin,
				BMI,
				DiabetesPedigreeFunction,
				Age,
				Outcome,
			); err != nil {
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
