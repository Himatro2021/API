package console

import (
	"himatro-api/internal/util"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "Himatro API",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		util.LogErr("ERROR", "command failed to execute", err.Error())
		log.Panic(err.Error())

		os.Exit(1)
	}
}
