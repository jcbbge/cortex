package cmd

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "cortex/internal/cortex"
    "cortex/internal/llm"
    "cortex/internal/spinner"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "gopkg.in/yaml.v3"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "cortex [question]",
    Short: "AI assistant in your terminal",
    Long: `Cortex is your AI assistant in the terminal.

Usage:
  cortex "what is docker compose?"    Ask a question
  ls                                  Use terminal normally`,
    Args: cobra.ArbitraryArgs,
    RunE: func(cmd *cobra.Command, args []string) error {
        // If no args, show simple welcome
        if len(args) == 0 {
            fmt.Println("Welcome to Cortex!")
            
            // Check if API key exists
            if viper.GetString("openai_api_key") == "" {
                fmt.Println("\nTo get started, you need an OpenAI API key.")
                fmt.Println("Enter your key now, or press Ctrl+C to exit.")
                fmt.Print("\nAPI Key: ")
                
                reader := bufio.NewReader(os.Stdin)
                key, err := reader.ReadString('\n')
                if err != nil {
                    return fmt.Errorf("failed to read API key: %w", err)
                }

                // Save the key
                home, err := os.UserHomeDir()
                if err != nil {
                    return fmt.Errorf("failed to get home directory: %w", err)
                }

                config := map[string]string{
                    "openai_api_key": strings.TrimSpace(key),
                }

                configBytes, err := yaml.Marshal(config)
                if err != nil {
                    return fmt.Errorf("failed to save config: %w", err)
                }

                configPath := filepath.Join(home, ".cortex.yaml")
                if err := os.WriteFile(configPath, configBytes, 0600); err != nil {
                    return fmt.Errorf("failed to write config: %w", err)
                }

                fmt.Println("\nGreat! Your API key is saved.")
                fmt.Println("Try asking a question:")
                fmt.Println("cortex \"what is docker?\"")
            } else {
                fmt.Println("\nAsk me anything:")
                fmt.Println("cortex \"your question here\"")
            }
            return nil
        }

        // Get API key
        apiKey := viper.GetString("openai_api_key")
        if apiKey == "" {
            fmt.Println("API key not found. Run 'cortex' to set it up.")
            return nil
        }

        // Initialize orchestrator
        orch, err := cortex.New(&llm.Config{APIKey: apiKey})
        if err != nil {
            return fmt.Errorf("initializing cortex: %w", err)
        }

        question := strings.Join(args, " ")
        responses, tracker, errs := orch.Process(context.Background(), question)

        s := spinner.New()
        s.Start()
        firstResponse := true

        for {
            select {
            case resp, ok := <-responses:
                if !ok {
                    responses = nil
                    fmt.Print("\n")
                    fmt.Print(tracker.FormatMetrics())
                    continue
                }
                if firstResponse {
                    s.Stop()
                    firstResponse = false
                    fmt.Print("\r\033[K") // Clear the entire line where spinner was
                }
                fmt.Print(resp)
            case err, ok := <-errs:
                if !ok {
                    errs = nil
                    continue
                }
                s.Stop()
                return fmt.Errorf("error: %w", err)
            default:
                if responses == nil && errs == nil {
                    return nil
                }
            }
        }
    },
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)
}

func initConfig() {
    home, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    // Look for config in home directory
    viper.AddConfigPath(home)
    viper.SetConfigType("yaml")
    viper.SetConfigName(".cortex")
    viper.ReadInConfig()
}