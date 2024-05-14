package cmd

import (
	"fmt"
	"github.com/pouspo/application-gateway/internal"
	"github.com/spf13/cobra"
)

var (
	DiagnosisCmd = &cobra.Command{
		Use:   "diagnosis",
		Short: "Manage Diagnosis",
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
