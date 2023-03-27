package resumes

import (
	"fmt"
	"strconv"
)

type Resumes []Resume

func (rr Resumes) PromptInput() string {
	var prompt string

	for _, r := range rr {
		if prompt != "" {
			prompt += "\n\n"
		}
		prompt += r.PromptInput()
	}

	return prompt
}

type Resume struct {
	ID   int
	Text string
}

func (r Resume) PromptInput() string {
	return "ResumeID: " + strconv.Itoa(r.ID) + "\n" + r.Text
}

type Candidate struct {
	Name  string
	Email string
}

type RecruiterMessageRequest struct {
	Candidate  Candidate
	ResumeText string
}

func (r RecruiterMessageRequest) PromptInput() string {
	return fmt.Sprintf("Candidate: %s\nEmail: %s\n Resume: %s", r.Candidate.Name, r.Candidate.Email, r.ResumeText)
}

type Email struct {
	To      string
	Subject string
	Body    string
}
