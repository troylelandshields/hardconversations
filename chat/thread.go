package chat

import (
	"context"

	gogpt "github.com/sashabaranov/go-openai"
	"github.com/troylelandshields/hardconversations/internal/tokens"
	"github.com/troylelandshields/hardconversations/logger"
	"github.com/troylelandshields/hardconversations/sources"
)

const baseSystemMessage = `You are an assistant interfacing with a machine, so answers must be given in the correct output with no surrounding punctuation. If you can't fulfill the request, respond with "Error: " and then a short explanation. %s

You can use the following information in responses:

###
`

const (
	roleUser      = "user"
	roleAssistant = "assistant"
	roleSystem    = "system"
)

type Thread struct {
	config Config
	ai     *gogpt.Client

	systemMessage       string
	systemMessageTokens int

	history           []gogpt.ChatCompletionMessage
	historyTokenCount int

	*sources.Manager
}

// NewThread returns a new Thread from a parent, inheriting the parent's history, sources, and config. The new Thread will have a new history. Different options can be applied to the new Thread.
func (t *Thread) NewThread(opt ...ConfigOption) *Thread {
	config := t.config
	for _, o := range opt {
		o(&config)
	}

	return &Thread{
		config: config,

		ai:                  t.ai,
		systemMessage:       t.systemMessage,
		systemMessageTokens: t.systemMessageTokens,

		history:           t.history,
		historyTokenCount: t.historyTokenCount,

		Manager: sources.NewFromParent(t.Manager),
	}
}

// Completely replaces existing history with the given history.
func (t *Thread) ReplaceHistory(history []gogpt.ChatCompletionMessage) {
	t.history = history
	for _, m := range history {
		t.historyTokenCount += tokens.MustCount(m.Content)
	}
}

func (t *Thread) ExecutePrompt(ctx context.Context, prompt string) (string, Metadata, error) {
	if t.systemMessageTokens == 0 {
		t.systemMessageTokens = tokens.MustCount(t.systemMessage)
	}

	// check if we need to drop any previous history
	if t.systemMessageTokens > t.config.MaxHistoryTokens {
		t.dropHistory(t.historyTokenCount - t.config.MaxHistoryTokens)
	}

	// push new user message to history
	t.pushHistory(roleUser, prompt)

	// find the source text information and append it to the system message
	contextInfoStr, usedSources, err := t.sourceText(
		ctx,
		t.config.MaxTotalTokens-
			(t.historyTokenCount+t.systemMessageTokens+t.config.MaxResponseTokens),
		prompt)
	if err != nil {
		return "", Metadata{}, err
	}
	systemMessage := t.systemMessage + contextInfoStr
	logger.Debugf("Sytem message: %s", systemMessage)

	messages := []gogpt.ChatCompletionMessage{
		{
			Role:    roleSystem,
			Content: systemMessage,
		},
	}
	messages = append(messages, t.history...)

	logger.Debugf("Sending question: %s", prompt)
	completionRequest := gogpt.ChatCompletionRequest{
		Model:       t.config.Model,
		Messages:    messages,
		MaxTokens:   t.config.MaxResponseTokens,
		Temperature: 0.0,
		TopP:        1.0,
	}

	resp, err := t.ai.CreateChatCompletion(ctx, completionRequest)
	if err != nil {
		return "", Metadata{}, err
	}

	responseText := resp.Choices[0].Message.Content

	logger.Debugf("Received answer: %s", responseText)
	t.pushHistory(roleAssistant, responseText)

	return responseText,
		Metadata{
			RawResponse:     resp,
			UsedTextSources: usedSources,
		},
		nil
}

func (t *Thread) sourceText(ctx context.Context, allowedTokens int, prompt string) (string, []sources.TextEmbedding, error) {
	sources, err := t.Manager.GetSourceText(ctx, t.config.UseEmbeddings, t.config.CosineSimilarityThreshold, allowedTokens, prompt, t.config.UserID)
	if err != nil {
		return "", nil, err
	}

	var contextualInfo string
	for _, s := range sources {
		if contextualInfo != "" {
			contextualInfo += "\n"
		}
		contextualInfo += s.Text
	}

	return contextualInfo, sources, nil
}

func (t *Thread) pushHistory(role, text string) {
	t.historyTokenCount += tokens.MustCount(text)

	t.history = append(t.history, gogpt.ChatCompletionMessage{
		Role:    role,
		Content: text,
	})
}

func (t *Thread) dropHistory(tokensToDrop int) {
	if tokensToDrop >= t.historyTokenCount {
		t.history = nil
		return
	}

	var droppedTokens int
	var dropToIdx int
	for i, msg := range t.history {
		droppedTokens += tokens.MustCount(msg.Content)
		if droppedTokens >= tokensToDrop {
			dropToIdx = i
			break
		}
	}

	t.history = t.history[dropToIdx:]
}

func (t *Thread) PurgeSources() {
	t.Manager = sources.New(t.ai)
}

type Metadata struct {
	RawResponse     gogpt.ChatCompletionResponse
	UsedTextSources []sources.TextEmbedding
}
