package moderatorai

import (
	"context"

	"github.com/troylelandshields/hardconversations/hardconversations"
)

const instruction = "Given the rules of a community and a piece of text, you are able to determine how likely it is that the text breaks the rules."

type Client struct {
	*hardconversations.Client
}

func NewClient(openAIKey string) *Client {
	c := &Client{
		Client: hardconversations.NewClient(openAIKey),
	}

	c.WithText(instruction)

	return c
}

type Thread struct {
	*hardconversations.Thread
}

func (c *Client) NewThread() *Thread {
	return &Thread{
		Thread: c.Client.NewThread(),
	}
}

const promptLikelihoodToBreakRules = `How likely it is that the given text breaks the rules?`

func (t *Thread) LikelihoodToBreakRules(ctx context.Context, input string) (result int, md hardconversations.Metadata, err error) {
	parsePrompt, err := hardconversations.ParseInstruction(result)
	if err != nil {
		return 0, hardconversations.Metadata{}, err
	}
	output, md, err := t.Thread.ExecutePrompt(ctx, parsePrompt+promptLikelihoodToBreakRules+"\n"+input)
	if err != nil {
		return 0, hardconversations.Metadata{}, err
	}

	err = hardconversations.Parse(output, &result)
	if err != nil {
		return 0, hardconversations.Metadata{}, err
	}

	return result, md, nil
}

const promptWhichRulesDoesItBreak = `Which rule numbers does the text break? (Answer must be a comma-separated list of integers)`

func (t *Thread) WhichRulesDoesItBreak(ctx context.Context) (result string, md hardconversations.Metadata, err error) {
	parsePrompt, err := hardconversations.ParseInstruction(result)
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}
	output, md, err := t.Thread.ExecutePrompt(ctx, parsePrompt+promptWhichRulesDoesItBreak)
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}

	err = hardconversations.Parse(output, &result)
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}

	return result, md, nil
}

const promptWhyDoesItBreakTheRules = `Why does it break the rules?`

func (t *Thread) WhyDoesItBreakTheRules(ctx context.Context) (result string, md hardconversations.Metadata, err error) {
	parsePrompt, err := hardconversations.ParseInstruction(result)
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}
	output, md, err := t.Thread.ExecutePrompt(ctx, parsePrompt+promptWhyDoesItBreakTheRules)
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}

	err = hardconversations.Parse(output, &result)
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}

	return result, md, nil
}
