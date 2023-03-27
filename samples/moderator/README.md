# Auto Moderator

You are a developer for a company called Inbreddit, an alt-right clone of a popular website. Communities within Inbreddit can have their own rules. You've been tasked with building a feature that will automatically detect if a post is likely to break the rules of a community and to flag which rules and why to the end user.

# moderator.yaml

```yaml
version: 1
conversations:
  - path: "./moderatorai"
    instruction: |
      Given the rules of a community and a piece of text, you are able to determine how likely it is that the text breaks the rules.
    # init_thread:
    #   - input: "This is a test"
    #     output: 0
    questions:
      - function_name: LikelihoodToBreakRules
        prompt: How likely is it that the text breaks the rules? (Answer must be an integer between 0 and 100)
        input: string
        output: int

      - function_name: WhichRulesDoesItBreak
        prompt: Which rule numbers does the text break? (Answer must be a comma-separated list of integers)
        output: "[]int"
        
      - function_name: WhyDoesItBreakTheRules
        prompt: Why does it break the rules?
        output: string
```