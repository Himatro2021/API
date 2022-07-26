package console

import (
	"log"
	"os"

	"github.com/Himatro2021/API/internal/helper"
	"github.com/spf13/cobra"
)

var runHash = &cobra.Command{
	Use:   "hash",
	Short: "hash a string",
	Long:  "hash a string input and output a hashed string",
	Run:   hash,
}

var runCheckHashMatch = &cobra.Command{
	Use:   "check-hash",
	Short: "check is hash match",
	Long:  "Check whether the plain param match with the hashed param",
	Run:   checkHash,
}

func init() {
	RootCmd.AddCommand(runHash)
	RootCmd.AddCommand(runCheckHashMatch)
}

func hash(_ *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("invalid command. Missing input argument (string)")
	}

	input := args[0]

	hashed, err := helper.HashString(input)
	if err != nil {
		log.Fatal("unexpected error happen", err.Error())
	}

	log.Println(hashed)
}

func checkHash(_ *cobra.Command, args []string) {
	if len(args) < 2 {
		log.Fatal("invalid command. Missing input argument (string)")
	}

	plain := args[0]
	hashed := args[1]

	if !helper.IsHashedStringMatch([]byte(plain), []byte(hashed)) {
		log.Println("Invalid!")
		log.Println("tips: try to surround your hashed args with singlequote e.g: 'your_hashed'")
		os.Exit(0)
	}

	log.Println("Match!")
	os.Exit(0)
}
