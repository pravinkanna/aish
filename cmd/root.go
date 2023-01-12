package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

type HttpRequest struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	MaxTokens   int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
}
type HttpResponse struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
}

type Choices struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     string `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aish/config.yaml)")
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
	if len(args) < 2 {
		return fmt.Errorf("please provide language and query to generate code snippet")
	}
	if len(args) > 16 {
		return fmt.Errorf("given prompt is too long")
	}
	lang := args[0]
	query := strings.Join(args[1:], " ")
	result := generateSnippet(lang, query)
	fmt.Println(result)
	return nil
}

func generateSnippet(lang string, query string) string {
	endpoint := `https://api.openai.com/v1/completions`
	apiKey := viper.GetString("api_key")
	model := viper.GetString("model")
	maxTokens := viper.GetInt("max_tokens")
	temperature := viper.GetInt("temperature")

	bearer := fmt.Sprintf("Bearer %s", apiKey)
	prompt := fmt.Sprintf("# %s\n# %s", lang, query)

	httpRequest := HttpRequest{model, prompt, maxTokens, temperature}
	postBody, _ := json.Marshal(httpRequest)
	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("POST", endpoint, requestBody)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error while reading the response bytes:", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(string([]byte(body)))
	}
	var response HttpResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		log.Fatal(err)
	}
	return response.Choices[0].Text
}
