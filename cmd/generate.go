package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code snippet",
	Long:  `Generate code snippet for the given input form ChatGPT`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: Please provide input to generate Code")
			return
		}
		if len(args) > 16 {
			fmt.Println("Error: Given input is too long")
			return
		}
		input := strings.Join(args[:], " ")
		fmt.Println("Your input is: ", input)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
