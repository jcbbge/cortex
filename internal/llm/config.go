package llm

import "errors"

// ModelType represents different model capabilities
type ModelType string

const (
    // OpenAI model names
    ModelReasoning ModelType = "gpt-4-turbo-preview" // Deep reasoning, complex analysis
    ModelCoding    ModelType = "gpt-4-turbo-preview" // High performance, coding tasks
    ModelLight     ModelType = "gpt-3.5-turbo"       // Quick, efficient responses

    // Model context limits
    ModelLightContextLimit    = 4096  // GPT-3.5-turbo context limit
    ModelReasoningContextLimit = 128000 // GPT-4-turbo context limit
)

// Config holds the LLM configuration
type Config struct {
    APIKey string `yaml:"openai_api_key"`
    Model  string // Model identifier, if not set uses appropriate default
}

func (c *Config) Validate() error {
    if c.APIKey == "" {
        return errors.New("OpenAI API key is required")
    }
    return nil
}

// GetContextLimit returns the context window size for the configured model
func (c *Config) GetContextLimit() int {
    switch c.Model {
    case string(ModelReasoning):
        return ModelReasoningContextLimit
    case string(ModelLight):
        return ModelLightContextLimit
    default:
        return ModelLightContextLimit // Default to smallest context window
    }
}

// DefaultConfig returns minimal required configuration
func DefaultConfig() *Config {
    return &Config{
        Model: string(ModelLight), // Default to lightweight model
    }
}