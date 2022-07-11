package console

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd :nodoc:
var RootCmd = &cobra.Command{
	Use: "Himatro API",
}

// Execute execute console command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)

		os.Exit(1)
	}
}
