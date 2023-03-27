# Bird Finder

> TODO: do images instead of text to keep with the spirit of the xkcd comic.

Given a page of text, count how many birds there are. Parse the details of each bird, and then generate a nice description.

See generated client usage in [main.go](https://github.com/troylelandshields/hardconversations/blob/main/samples/birdfinder/main.go).

# birdfinder.yaml

```yaml
version: 1
conversations:
  - path: "./birdai"
    instruction: |
      Given a piece of text, you are able to determine how many birds are mentioned in the text and describe each bird.
    questions:
      - function_name: CountBirds
        prompt: How many birds are mentioned in the text?
        input: string
        output: int
        
      - function_name: ParseBird
        prompt: Can you parse the details of each bird? 
        output: "[]github.com/troylelandshields/hardconversations/samples/birdfinder/bird.Bird"

      - function_name: DescribeBird 
        input: github.com/troylelandshields/hardconversations/samples/birdfinder/bird.Bird
        output: string
        prompt: Describe the bird with the given properties and add a fun fact (make it up if you have to)
```

## Usage
```go
birdCount, _, err := t.CountBirds(ctx, str)
if err != nil {
  fmt.Println(err)
  os.Exit(1)
}
```
