package birdai

import (
	"context"
	"encoding/json"

	"github.com/troylelandshields/hardconversations/hardconversations"
	"github.com/troylelandshields/hardconversations/samples/birdfinder/bird"
)

const (
	preamble = ""
)

type Client struct {
	*hardconversations.Client
}

func NewClient(openAIKey string) *Client {
	return &Client{
		Client: hardconversations.NewClient(openAIKey),
	}
}

type Thread struct {
	*hardconversations.Thread
}

func (c *Client) NewThread() *Thread {
	return &Thread{
		Thread: c.Client.NewThread(),
	}
}

const promptIsABird = `Is the text about a bird?`

func (t *Thread) IsABird(ctx context.Context) (result bool, md hardconversations.Metadata, err error) {
	parsePrompt, err := hardconversations.ParseInstruction(result)
	if err != nil {
		return false, hardconversations.Metadata{}, err
	}
	output, md, err := t.Thread.ExecutePrompt(ctx, parsePrompt+promptIsABird)
	if err != nil {
		return false, hardconversations.Metadata{}, err
	}

	err = hardconversations.Parse(output, &result)
	if err != nil {
		return false, hardconversations.Metadata{}, err
	}

	return result, md, nil
}

const promptParseBird = `Parse the details of the bird?`

func (t *Thread) ParseBird(ctx context.Context) (bird.Bird, hardconversations.Metadata, error) {
	parsePrompt, err := hardconversations.ParseInstruction(bird.Bird{})
	if err != nil {
		return bird.Bird{}, hardconversations.Metadata{}, err
	}

	output, md, err := t.Thread.ExecutePrompt(ctx, parsePrompt+promptParseBird)
	if err != nil {
		return bird.Bird{}, hardconversations.Metadata{}, err
	}

	var result bird.Bird
	err = hardconversations.Parse(output, &result)
	if err != nil {
		return bird.Bird{}, hardconversations.Metadata{}, err
	}

	return result, md, nil
}

const promptDescibeTheBird = `Describe the bird with the given properties and add a fun fact (make it up if you have to)`

func (t *Thread) DescribeBird(ctx context.Context, input bird.Bird) (result string, md hardconversations.Metadata, err error) {
	parsePrompt, err := hardconversations.ParseInstruction(result)
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}

	inputString, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}

	output, md, err := t.Thread.ExecutePrompt(ctx, parsePrompt+promptDescibeTheBird+"\n"+string(inputString))
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}

	err = hardconversations.Parse(output, &result)
	if err != nil {
		return "", hardconversations.Metadata{}, err
	}

	return result, md, nil
}
