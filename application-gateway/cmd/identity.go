package identity

import (
	"fmt"
	"github.com/pouspo/application-gateway/internal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
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
)
