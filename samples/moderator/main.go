package main

import (
	"context"
	"fmt"
	"os"

	"github.com/troylelandshields/hardconversations/samples/moderator/moderatorai"
)

const theRules = `The rules are simple:

1. Don't talk about Fight Club.
2. Don't talk about Fight Club.
3. If someone says 'stop' or goes limp, taps out the fight is over.
`

func main() {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		fmt.Println("OPENAI_API_KEY env variable is required")
		os.Exit(1)
	}

	autoModClient := moderatorai.NewClient(openAIKey)
	ctx := context.Background()

	autoModClient.WithText(theRules)

	t := autoModClient.NewThread()

	likelihood, _, err := t.LikelihoodToBreakRules(ctx, `"The thing that I love about Fight Club is getting out my aggression and posting pictures online."`)
	// likelihood, _, err := t.LikelihoodToBreakRules(ctx, `"I like to hang out with my friends and do nothing in particular at all."`)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Println("likelihood:", likelihood)

	if likelihood < 50 {
		fmt.Println("no rule breaking here")
		return
	}

	rules, _, err := t.WhichRulesDoesItBreak(ctx)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	reason, _, err := t.WhyDoesItBreakTheRules(ctx)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Println("rules:", rules, reason)
}
