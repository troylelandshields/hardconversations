package main

import (
	"context"
	"fmt"
	"os"

	"github.com/troylelandshields/hardconversations/chat"
	"github.com/troylelandshields/hardconversations/logger"
	"github.com/troylelandshields/hardconversations/samples/recruiter/autorecruiter"
	"github.com/troylelandshields/hardconversations/samples/recruiter/resumes"
	"github.com/troylelandshields/hardconversations/sources"
)

const (
	developerJobDescription       = `Senior Backend Engineer, with at least 5 years of industry experience. Prefer professional experience with Go, but will consider other languages. Must have experience with microservices, and be able to work in a fast-paced environment. Salary range of $250,000 - $275,000. A degree is preferred, but not required.`
	graphicDesignerJobDescription = `Senior Graphic Designer, preferred to have at least 5 years of industry experience as a designer. Must have experience with Adobe Creative Suite, and be able to work in a fast-paced environment. Salary range of $150,000 - $175,000. Any degree is required.`
	jobState                      = "UT"
)

func main() {
	logger.SetLogLevel(logger.LevelInfo)

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		fmt.Println("OPENAI_API_KEY env variable is required")
		os.Exit(1)
	}

	aiRecruiter := autorecruiter.NewClient(openAIKey, chat.WithUseEmbeddings(true), chat.WithCosineSimilarityThreshold(0.7))

	// returns all in-state resumes
	aiRecruiter.AddSourceTextProvider(resumes.ResumeProvider{})

	// out of state resumes are less likely to be a good fit, so we'll provide them but with a lower weight.
	aiRecruiter.AddSourceTextProvider(resumes.OutOfStateResumeProvider{}, sources.WithWeight(0.95))

	ctx := context.Background()
	ctx = context.WithValue(ctx, "job_state", jobState)

	t := aiRecruiter.NewThread()

	rankedResumeIDs, _, err := t.RankResumes(ctx, developerJobDescription)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	// don't need resume sources anymore so we can remove it.
	t.PurgeSources()

	for _, id := range rankedResumeIDs {
		resume, err := resumes.LookupResume(id)
		if err != nil {
			fmt.Println("error getting resume, will try next one", err)
			continue
		}

		candidate, _, err := t.GetCandidateInfo(ctx, resume.Text)
		if err != nil {
			fmt.Println("error getting cadnidate info, will try next one", err)
			continue
		}

		msg, _, err := t.GenerateRecruiterMessage(ctx, resumes.RecruiterMessageRequest{
			Candidate:  candidate,
			ResumeText: resume.Text,
		})
		if err != nil {
			fmt.Println("error generating recruiter message, will try next one", err)
			continue
		}

		fmt.Printf("%+v\n\n", msg)
	}
}
