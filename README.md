# Hard Conversations
#### Generate a statically typed client from YAML to interface with ChatGPT for whatever problems you're trying to solve.
---

## Introduction

ChatGPT technology is useful for much more than a user-to-bot chat assistant. Developers can utilize it as a general-purpose, "basic reasoning" API to enhance the types of problems that are easily solved with software, without having to gather and categorize data to train machine-learning models. Hard Conversations comes with a CLI tool to generate a statically typed client to make computer-to-bot interfacing much easier. Jump to the example below to see how works.

### Quick Example

You have a string and you want to count how many US Presidents are mentioned in it.

YAML -> generate -> Code

Define your question in YAML.

```yaml
version: 1
conversations:
  - path: "./presidents"
    instruction: |
      You can count how many US Presidents are listed in some text.
    questions:
      - function_name: CountPresidents
        prompt: How many presidents are mentioned?
        input: "string"
        output: "int"
```

Generate a client with `hardc generate`.

Use your statically typed client to easily build an application.

```go
	client := presidents.NewClient(openAIKey)
	countOfPresidents, _, _ := client.CountPresidents()

	if countOfPresidents > 3 {
		fmt.Println("That's a lot of presidents")
		return
	}
```

[Bigger Example Below](https://github.com/troylelandshields/hardconversations/blob/main/README.md#example)

### Quickstart

#### Installation

```bash
go install github.com/troylelandshields/hardconversations/cmd/hardc
```

#### Usage

```bash
hardc generate # defaults to hardc.yaml
hardc generate -f path/to/file.yaml
```

# Background

## Soft Inputs

This comic is the perfect example of a problem that used to be difficult but is now trivially solved by sprinkling some ChatGPT onto your program with the help of HardConversations.

![](https://imgs.xkcd.com/comics/tasks.png)

Just a year ago you would have thought this comic was a hilarious and biting commentary on the ridiculous requests that Product Managers and the Non-Technicals harry us with. 

> Detect if the image has a BIRD?! Ha! It's virtually impossible, wait until /r/ProgrammerHumor gets a load of this!

It's 2023 though, and there's a new sheriff in town.

## No More Excuses

Solving fuzzy or "soft" problems like bird-detection has been pretty attainable for a while now with machine-learning, but the upfront cost of collecting enough data, training a model, and integrating into your application might still have been too much to take on for a small team or an experimental feature. So while techincally the problem could have been solved, it just may not have been worth the effort.

ChatBots will imminently become ubiquitous in most software to help users, but if that's the only application of this technology we can come up with we're missing some huge opportunities in my opinion. What if we treat ChatGPT as a general-purpose AI that we can program against to quickly and easily build out features like a bird-detector? What if a feature that would have taken a research team 5 years to figure out in the past could be completed in a couple hours by an unpaid intern with a can-do attitude? ChatGPT can help you solve problems faster and more easily than traditional software development techniques, allowing you to quickly prototype and test new features.

Don't worry about fuzzy or soft problems anymore, because with ChatGPT your programs can now have *HardConversations*.

## Hard Conversations

This tool takes a YAML file as input and generates a client to converse with OpenAI's ChatGPT. In your YAML file, you will list the specific questions you want to be able to ask with the expected inputs and outputs, and HardConversations will generate a client that you can then easily use in your program as you see fit.

# Example

Let's say you work for a recruiting agency. Your company has a database full of resumes of potential job candidates. Clients come to you with job descriptions. You are tasked with building an automated tool that does all of the following:

1. Find the 3 best candidates for a job description. Candidates can live anywhere, but we should prefer candidates that live in the same state as the job.
2. Parse the candidates contact info from their resume.
3. Automatically send an email to the candidates with a description of the job and an invitation to apply.

Let's see how HardConversations can help us solve these problems.

We can define the details of the "conversation" that our application wants to have in a YAML file to make answering these questions easy. We need to give some instructions about the conversation to give some context, and then we define the various types of questions our application can "ask." We will "ask" these questions by calling functions. These questions are going to be directly related to the 3 requirements above.

We also need to define what the input and output should be for each of these questions. For example, if we want to parse candidate info from a resume, we can have HardConversations make sure the answer is returned from ChatGPT as a Go struct called `Candidate`, which has fields for `Name` and `Email`.

Putting this all together in a YAML file looks like this:
 
```yaml
version: 1
conversations:
  - path: "./autorecruiter"
    instruction: |
      Given a list of resumes, you are able to determine which ones are the best fit for the job description.
    questions:
      - function_name: RankResumes
        prompt: Return just the IDs of between 1 and 3 resumes in a comma-separated list, ranked from best to worst fit for the job description. Do not include resumes that are not a good fit.
        input: "string"
        output: "[]int"

      - function_name: GetCandidateInfo
        prompt: "Return the candidate info from the resume"
        input: string
        output: github.com/troylelandshields/hardconversations/samples/recruiter/resumes.Candidate

      - function_name: GenerateRecruiterMessage
        prompt: Generate a message to send to the candidate about the job; mention what you like about their resume and why you think they would be a good fit for the job.
        input: "github.com/troylelandshields/hardconversations/samples/recruiter/resumes.RecruiterMessageRequest"
        output: github.com/troylelandshields/hardconversations/samples/recruiter/resumes.Email
```

Now that we've defined this "conversation", we want to be able to write an application that can use this functionality. We use the `hardc` CLI to generate libraries from this YAML file by executing `hardc generate -f path/to/file.yaml`.

Now we have an auto-generated, statically-typed client that we can create like this:

```go
aiRecruiter := autorecruiter.NewClient(openAIKey)
```
However, that's not all we need for this to work. We need to provide the resumes as a data-source that ChatGPT can utilize to find the best candidates.

To do that, we can add a "source provider" that HardConversations can use to add more contextual information to the conversation with ChatGPT, as needed.

```go
// resume.ResumeProvider will return a list of resumes from the same state as the job.
aiRecruiter.AddSourceTextProvider(db.ResumeProvider{})

// we also want to provide out-of-state resumes, but with slightly less preference, so we'll weight them a little less.
aiRecruiter.AddSourceTextProvider(db.OutOfStateResumeProvider{}, sources.WithWeight(0.95))
```

Now, we can match resumes with jobs. Let's say we wanted to do it in a web-service. It would look something like this rough outline:

```go
func HandleNewJob(w http.ResponseWriter, r *http.Request) {
	// get new job details from the request; description and state
	var jobDetails Job	
	json.NewDecoder(r.Body).Decode(&jobDetails)

	// add the job state to the context so it can be used by the resume provider
	ctx := context.WithValue(r.Context(), jobStateKey, job.State)
	
	// create a new thread for this "conversation"
	thread := aiRecruiter.NewThread()

	// ask ChatGPT to rank the best fitting resumes; the provided sources will be used as contextual info
	resumeIDs, _, _ := thread.RankResumes(ctx, jobDetails.Description)

	for _, id := range resumeIDs {
		// get the resume from the database
		resume, _ := db.LookupResume(id)

		// ask ChatGPT to parse the Candidate details from the unstructured text of the resume
		candidate, _, _ := thread.GetCandidateInfo(ctx, resume.Text)

		// ask ChatGPT to generate a personalized message that we can send to the candidate
		personalizedMessage, _, _ := thread.GenerateRecruiterMessage(ctx, RecruiterMessageRequest{
			Candidate: candidate,
			ResumeText: resume.Text,
		})

		// ... send email message
	}
}
```

Solving any single one of these problems in software may have been pretty difficult. Interfacing with ChatGPT in a statically typed way through the client generated by HardConversations, however, we were able to get a working version up-and-running very quickly, and now we can iterate from here to continually measure and improve its success metrics.

## More Samples

* [Automoderator](https://github.com/troylelandshields/hardconversations/tree/main/samples/moderator)
* [Bird-finder](https://github.com/troylelandshields/hardconversations/tree/main/samples/birdfinder)
* [AI Recruiter](https://github.com/troylelandshields/hardconversations/tree/main/samples/recruiter)
* TODO: HotDog/NotHotDog (once I have access to the ChatGPT4 and can use images as an input)

## Acknowledgements	

Thanks to [go-openai](https://github.com/sashabaranov/go-openai)

Inspired by SQLc.

