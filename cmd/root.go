package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aish",
	Short: "An AI based code generator",
	Long: `AIsh is an AI-powered command line tool that generates code snippets, 
helping developers to save time and increase productivity
This project is powered by OpenAI's GPT-3`,
	RunE: runCmd,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aish)")
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home + `/.aish`)
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func runCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide an input to generate code snippet")
	}
	if len(args) > 16 {
		return fmt.Errorf("given input is too long")
	}
	input := strings.Join(args[:], " ")
	result := generateSnippet(input)
	fmt.Println(result)
	return nil
}

func generateSnippet(input string) string {
	return "dummy \noutput"
}
