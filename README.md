# Hard Conversations

### Sprinkle some ChatGPT onto your program
---
####Generate a statically typed client from YAML to interface with ChatGPT for whatever problems you're trying to solve.
---

## Quickstart

### Installation

```bash
go install github.com/troylelandshields/hardconversations/cmd/hardc
```

### Usage

```bash
hardc generate # defaults to hardc.yaml
hardc generate -f path/to/file.yaml
```

## Soft Inputs

Read this comic:

![](https://imgs.xkcd.com/comics/tasks.png)

Just a year ago you would have thought this comic was a hilarious and biting commentary on the ridiculous requests that Product Managers and the Non-Technicals harry us with. 

> Detect if the image has a BIRD?! Ha! It's virtually impossible, wait until /r/ProgrammerHumor gets a load of this!

It's 2023 though, and there's a new sheriff in town.

## No More Excuses

Solving fuzzy or "soft" problems like bird-detection has been pretty attainable for a while now with machine-learning, but the upfront cost of collecting enough data, training a model, and integrating into your application might still have been too much to take on for a small team or an experimental feature. So while techincally the problem could have been solved, it just may not have been worth the effort.

ChatBots will imminently become ubiquitous in most software to help users, but if that's the only application of this technology we can come up with we're missing some huge opportunities in my opinion. What if we treat ChatGPT as a general-purpose AI that we can program against to quickly and easily build out features like a bird-detector? What if a feature that would have taken a research team 5 years to figure out in the past could be completed in a couple hours by an unpaid intern with a can-do attitude? What if we stop worrying about what's possible or not possible, and instead we just try it out and see how it goes?

Don't worry about fuzzy or soft problems anymore, because with ChatGPT your programs can now have *HardConversations*.

## Hard Conversations

This tool takes a YAML file as input and generates a client to converse with OpenAI's ChatGPT. In your YAML file, you will list the specific questions you want to be able to ask with the expected inputs and outputs, and HardConversations will generate a client that you can then easily use in your program as you see fit.

## Example

For example:
 
```yaml
version: 1
packages:
  - path: "./moderatorai"
    instruction: |
      Given the rules of a community and a piece of text, you are able to determine how likely it is that the text breaks the rules.
    questions:
      - functionName: LikelihoodToBreakRules
        prompt: How likely is it that the text breaks the rules? (Answer must be an integer between 0 and 100)
        input: string
        output: int
      - functionName: WhichRulesDoesItBreak # to flag specific rules in the UI
        prompt: Which rule numbers does the text break?
        output: []int
      - functionName: WhyDoesItBreakTheRules # to show an explanation to users
        prompt: Why does it break the rules?
        output: string
```

Generates a client that you can use like this:

```go
	autoModClient := moderatorai.NewClient(openAIKey)
	autoModClient.WithText(theRules)
	t := autoModClient.NewThread()

	likelihood, _, err := t.LikelihoodToBreakRules(ctx, `"The thing that I love about Fight Club is getting out my aggression and posting pictures online."`)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

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
```

## More Samples

* [Inbreddit Automoderator](https://github.com/troylelandshields/hardconversations/tree/main/samples/moderator)
* [AI Recruiter](https://github.com/troylelandshields/hardconversations/tree/main/samples/recruiter)
* [Bird-finder](https://github.com/troylelandshields/hardconversations/tree/main/samples/birdfinder)
* TODO: HotDog/NotHotDog (once I have access to the ChatGPT4 and can use images as an input)

## Acknowledgements	

Inspired by SQLc.

