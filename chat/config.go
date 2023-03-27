package chat

import (
	gogpt "github.com/sashabaranov/go-openai"
)

type Config struct {
	MaxTotalTokens    int // defaults to 4000
	MaxResponseTokens int // defaults to 300
	MaxHistoryTokens  int // defaults to 1000

	// ChatRequest options
	Model       string  // defaults to gogpt.GPT3Dot5Turbo
	Temperature float64 // defaults to 0, max 2
	UserID      string  // defaults to ""

	// TODO: support
	UseEmbeddings             bool    // defaults to false
	CosineSimilarityThreshold float64 // defaults to 0.7, must be between 0 and 1.
	// MaxTokensChunkSize        int // TODO: figure out chunking
}

// NewConfig returns a new Config with default values.
func NewConfig(opt ...ConfigOption) Config {
	defaults := &Config{
		MaxTotalTokens:    4000,
		MaxResponseTokens: 300,
		MaxHistoryTokens:  1000,

		Model:       gogpt.GPT3Dot5Turbo,
		Temperature: 0,
		UserID:      "",

		UseEmbeddings:             false,
		CosineSimilarityThreshold: 0.7,
	}

	for _, o := range opt {
		o(defaults)
	}

	return *defaults
}

type ConfigOption func(*Config)

func WithMaxTotalTokens(maxTotalTokens int) ConfigOption {
	return func(c *Config) {
		c.MaxTotalTokens = maxTotalTokens
	}
}

func WithMaxResponseTokens(maxResponseTokens int) ConfigOption {
	return func(c *Config) {
		c.MaxResponseTokens = maxResponseTokens
	}
}

func WithMaxHistoryTokens(maxHistoryTokens int) ConfigOption {
	return func(c *Config) {
		c.MaxHistoryTokens = maxHistoryTokens
	}
}

func WithModel(model string) ConfigOption {
	return func(c *Config) {
		c.Model = model
	}
}

func WithTemperature(temperature float64) ConfigOption {
	return func(c *Config) {
		c.Temperature = temperature
	}
}

func WithUserID(userID string) ConfigOption {
	return func(c *Config) {
		c.UserID = userID
	}
}

func WithUseEmbeddings(useEmbeddings bool) ConfigOption {
	return func(c *Config) {
		c.UseEmbeddings = useEmbeddings
	}
}

// WithCosineSimilarityThreshold changes the minimum value required to include a source in the chat request.
// The cosine similarity between a source and the prompt (times source weight) must meet this minimum or it will be ignored.
// Only used if UseEmbeddings is true.
func WithCosineSimilarityThreshold(cosineSimilarityThreshold float64) ConfigOption {
	return func(c *Config) {
		c.CosineSimilarityThreshold = cosineSimilarityThreshold
	}
}

// TODO:
// func WithMaxTokensChunkSize(maxTokensChunkSize int) ConfigOption {
// 	return func(c *Config) {
// 		c.MaxTokensChunkSize = maxTokensChunkSize
// 	}
// }
