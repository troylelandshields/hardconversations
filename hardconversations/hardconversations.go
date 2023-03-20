package hardconversations

import (
	"context"
	"fmt"
	"strings"

	gogpt "github.com/sashabaranov/go-openai"
)

const systemMessage = `You are an assistant interfacing with a machine, so answers must be given in the correct output with no surrounding punctuation.

You will answer questions from the following information: ###
%s
###
`

type Client struct {
	ai *gogpt.Client

	*Thread // starting thread
}

func NewClient(openAIKey string) *Client {
	openAIClient := gogpt.NewClient(openAIKey)
	return &Client{
		ai:     openAIClient,
		Thread: &Thread{},
	}
}

type Thread struct {
	ai            *gogpt.Client
	systemMessage string
	history       []gogpt.ChatCompletionMessage

	textualInfo []string
}

func (c *Client) NewThread() *Thread {
	return &Thread{
		ai:            c.ai,
		systemMessage: systemMessage,
		history:       c.Thread.history,
		textualInfo:   c.Thread.textualInfo,
	}
}

func (t *Thread) WithText(text string) {
	t.textualInfo = append(t.textualInfo, text)
}

func (t *Thread) ExecutePrompt(ctx context.Context, prompt string) (string, Metadata, error) {
	t.history = append(t.history, gogpt.ChatCompletionMessage{
		Role:    "user",
		Content: prompt,
	})

	systemMessage := fmt.Sprintf(t.systemMessage, strings.Join(t.textualInfo, "\n"))

	messages := []gogpt.ChatCompletionMessage{
		{
			Role:    "system",
			Content: systemMessage,
		},
	}
	messages = append(messages, t.history...)

	completionRequest := gogpt.ChatCompletionRequest{
		Model:       gogpt.GPT3Dot5Turbo,
		Messages:    messages,
		MaxTokens:   300,
		Temperature: 0.0,
		TopP:        1.0,
	}

	fmt.Println("completionRequest", completionRequest.Messages)
	resp, err := t.ai.CreateChatCompletion(ctx, completionRequest)
	if err != nil {
		return "", Metadata{}, err
	}

	responseText := resp.Choices[0].Message.Content

	t.history = append(t.history, gogpt.ChatCompletionMessage{
		Role:    "assistant",
		Content: responseText,
	})

	return responseText,
		Metadata{
			RawResponse: resp,
		},
		nil
}

type Metadata struct {
	RawResponse gogpt.ChatCompletionResponse
}
