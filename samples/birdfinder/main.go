package main

import (
	"context"
	"fmt"
	"os"

	"github.com/troylelandshields/hardconversations/samples/birdfinder/birdai"
)

func main() {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		fmt.Println("OPENAI_API_KEY env variable is required")
		os.Exit(1)
	}

	birdFinderClient := birdai.NewClient(openAIKey)
	ctx := context.Background()

	t := birdFinderClient.NewThread()

	t.WithText("An african swallow can carry a coconut.")
	// t.WithText("The bird is trapped in a cage.")

	// t.WithEmbeddedText("")
	// t.WithImage("")
	// t.WithTextSource() // takes an interface that returns sorted relevant text from somewhere
	// t.WithImageSource() ?

	isBird, md, err := t.IsABird(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(md)

	if !isBird {
		fmt.Println("no bird here")
		return
	}

	bird, _, err := t.ParseBird(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", bird)

	desc, _, err := t.DescribeBird(ctx, bird)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(desc)
}
