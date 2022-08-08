package console

import (
	"github.com/Himatro2021/API/internal/helper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "encrypt a string",
	Long:  "Encrypt a string using encryption standars implemented in this app",
	Run:   encrypt,
}

func init() {
	RootCmd.AddCommand(encryptCmd)
}

func encrypt(cmd *cobra.Command, args []string) {
	input := args[0]
	if input == "" {
		logrus.Fatal("input not found")
	}

	encrypted, err := helper.Cryptor().Encrypt(input)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("result: ", encrypted)
}
