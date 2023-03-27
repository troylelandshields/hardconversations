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

	// str := "I still live with my mom and dad. I'm 35 years old. They were probably expecting an empty nest at this point, but I need to save money on rent."
	// str := "An african swallow can carry a coconut."
	// str := "The bird is trapped in a cage."
	// str := "The Andean Condor and the American Bald Eagle are both birds of prey, but they live in different parts of the world. The former is native to South America, whereas the latter is native to North America. The Andean Condor is the largest flying bird in the world, with a wingspan of up to 3.2 meters (10.5 feet). The American Bald Eagle is the national bird of the United States of America."
	// t.WithImageSource() ?

	birdCount, _, err := t.CountBirds(ctx, str)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if birdCount == 0 {
		fmt.Println("no bird here")
		return
	}

	birds, _, err := t.ParseBird(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, bird := range birds {
		fmt.Printf("%d. %+v\n", i+1, bird)

		// desc, _, err := t.DescribeBird(ctx, bird)
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }
		// fmt.Println(desc)
	}
}
