package console

import (
	"bufio"
	"fmt"
	"himatro-api/internal/auth"
	"himatro-api/internal/util"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var encryptorCmd = &cobra.Command{
	Use:   "encryptor",
	Short: "Encrypt your input",
	Long:  "If you want to initialize admin password, u can use this.",
	Run:   encryptor,
}

func init() {
	RootCmd.AddCommand(encryptorCmd)
}

func encryptor(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Input your text to be encrypted (enter to proceed)")
	fmt.Print("input: ")

	text, err := reader.ReadString(byte('\n'))

	if err != nil {
		util.LogErr("ERROR", "command encryptor is fail", err.Error())
		fmt.Println("command failed")
		log.Panic(err.Error())
	}

	encrypted, err := auth.Encrypt(text)

	if err != nil {
		util.LogErr("ERROR", "command encryptor is fail", err.Error())
		fmt.Println("command failed")
		log.Panic(err.Error())
	}

	fmt.Println(encrypted)
}
