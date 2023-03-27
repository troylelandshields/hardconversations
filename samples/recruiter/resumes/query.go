package resumes

import (
	"context"
	"fmt"
)

var (
	db = map[string]Resumes{
		"UT": {
			{
				ID:   3,
				Text: seniorDeveloperResume,
			},
			{
				ID:   4,
				Text: midLevelDeveloperResume,
			},
		},
		"CA": {
			{
				ID:   1,
				Text: recentGradResume,
			},
			{
				ID:   2,
				Text: graphicDesignerResume,
			},
		},
	}
)

func LookupResumes(state string) (Resumes, error) {
	resumes, ok := db[state]
	if !ok {
		return nil, nil
	}

	return resumes, nil
}

func LookupResume(id int) (Resume, error) {
	for _, state := range allStates {
		resumes, err := LookupResumes(state)
		if err != nil {
			return Resume{}, err
		}

		for _, resume := range resumes {
			if resume.ID == id {
				return resume, nil
			}
		}
	}

	return Resume{}, fmt.Errorf("no resume with id %s", id)
}

type ResumeProvider struct {
}

func (r ResumeProvider) Sources(ctx context.Context, prompt string) ([]string, error) {
	state, ok := ctx.Value("job_state").(string)
	if !ok {
		return nil, fmt.Errorf("no job state in context")
	}

	resumes, err := LookupResumes(state)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, resume := range resumes {
		result = append(result, resume.PromptInput())
	}

	return result, nil
}

type OutOfStateResumeProvider struct {
}

func (r OutOfStateResumeProvider) Sources(ctx context.Context, prompt string) ([]string, error) {
	jobState, ok := ctx.Value("job_state").(string)
	if !ok {
		return nil, fmt.Errorf("no job state in context")
	}

	var result []string
	for _, state := range allStates {
		if jobState == state {
			continue
		}

		resumes, err := LookupResumes(state)
		if err != nil {
			return nil, err
		}

		for _, resume := range resumes {
			result = append(result, resume.PromptInput())
		}
	}

	return result, nil
}

var (
	allStates = []string{
		"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA", "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC", "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY",
	}
)

const (
	recentGradResume = `Janet Dough
	j@dough.com
	
	EDUCATION
	Bachelor of Science in Computer Science, XYZ University (May 2022)
	
	EXPERIENCE
	Software Engineering Intern, ABC Corporation (May-August 2021)
	
	Worked on a team responsible for developing and maintaining a web application used by millions of users
	Collaborated with senior engineers to identify and fix bugs in the codebase
	Contributed to the development of new features using React and Node.js
	Gained experience working in an agile development environment

	SKILLS
	Programming languages: Java, Python, JavaScript
	Web development: React, Node.js, HTML, CSS
	Database management: SQL, MongoDB
	Other: Git, Agile methodology
	
	PROJECTS
	
	Developed a simple web application using React and Node.js to manage and display a personal book collection
	Implemented a basic machine learning model in Python to classify images of handwritten digits

	ACTIVITIES
	
	Member of XYZ University Computer Science Club
	Volunteered at local coding workshops for high school students`

	graphicDesignerResume = `Bobo Roberts
	boborobo@email.com
	
	PROFILE
	Graphic designer with 10+ years of experience creating impactful designs for a variety of industries, including technology, fashion, and entertainment. Skilled in Adobe Creative Suite, with a strong eye for detail and a passion for creating visually stunning designs.
	
	EXPERIENCE
	Senior Graphic Designer, Apple Inc. (2017-2022)
	
	Led the design and development of marketing materials for new product launches, resulting in a 20% increase in sales
	Collaborated with cross-functional teams to create and maintain Apple's brand guidelines
	Designed visual assets for Apple's website and social media channels, including graphics, animations, and videos
	Mentored junior designers, providing guidance on design principles and best practices

	Graphic Designer, XYZ Company (2014-2017)
	
	Worked on a team responsible for creating advertising campaigns for a variety of clients in the fashion industry
	Designed print and digital materials, including billboards, magazine ads, and email newsletters
	Conducted research to understand client needs and target audience, resulting in campaigns that resonated with consumers
	Presented design concepts to clients and incorporated feedback to deliver final designs that exceeded expectations
	
	Graphic Designer, ABC Agency (2012-2014)
	
	Created visual designs for websites and mobile applications for clients in the entertainment industry
	Collaborated with UX designers and developers to ensure designs were both visually stunning and user-friendly
	Designed logos, brand identities, and marketing materials for new and existing clients
	Contributed to the development of the agency's internal design standards and processes

	SKILLS
	Adobe Creative Suite (Photoshop, Illustrator, InDesign, After Effects)
	Graphic design principles and best practices
	Typography and layout design
	Branding and identity design
	User experience design
	
	AWARDS
	
	Winner, AIGA 50 Books/50 Covers Competition (2019)
	Finalist, Communication Arts Interactive Design Competition (2016)
	Honorable Mention, Print Magazine Regional Design Awards (2014)

	EDUCATION
	Bachelor of Fine Arts in Graphic Design, XYZ University (2012)`

	midLevelDeveloperResume = `David Lee
	davidlee@email.com
	
	SUMMARY
	Mid-level engineer with 3 years of experience in software development and a strong foundation in computer science concepts. Skilled in programming languages such as Java, Python, and C++, as well as experience with various software development tools and technologies. Currently pursuing a degree in Computer Science at ABC University.
	
	EXPERIENCE
	Software Engineer, XYZ Corporation (2019-Present)
	
	Designed and developed software features for a large-scale enterprise web application using Java and Spring Framework
	Collaborated with cross-functional teams, including product managers and QA engineers, to deliver high-quality software on time
	Participated in code reviews, providing feedback and suggestions to improve code quality and maintainability
	Actively contributed to the development of the company's internal software development standards and best practices
	Software Engineering Intern, DEF Company (2018)
	
	Developed software features for an e-commerce web application using Python and Django Framework
	Worked with senior engineers to identify and fix bugs in the codebase
	Contributed to the development of new features using React and Node.js
	Gained experience working in an agile development environment

	SKILLS
	Programming languages: Java, Python, C++
	Web development: Spring Framework, Django Framework, React, Node.js, HTML, CSS
	Database management: SQL
	Other: Git, Agile methodology, Scrum, JIRA
	
	EDUCATION
	Bachelor of Science in Computer Science (In Progress), ABC University
	
	Expected graduation: May 2024
	Relevant coursework: Data Structures and Algorithms, Object-Oriented Programming, Computer Networks, Operating Systems
	Associate of Science in Computer Science, XYZ Community College (2018)
	
	Relevant coursework: Java Programming, Database Systems, Web Development

	ACTIVITIES
	
	Volunteer, local hackathons and coding workshops for underrepresented groups
	Member, XYZ University Computer Science Club`

	seniorDeveloperResume = `Mark Johnson
	markjohnson@email.com
	
	PROFILE
	Senior backend developer with 7+ years of experience developing scalable and reliable systems. Proficient in programming languages such as Go and Elixir, with expertise in developing microservices and APIs. Experienced in leading teams and mentoring junior developers.
	
	EXPERIENCE
	Senior Backend Developer, XYZ Inc. (2019-Present)
	
	Led the development of a microservice architecture using Go and Elixir, resulting in a 50% increase in system performance
	Collaborated with cross-functional teams to design and develop APIs for a mobile banking application
	Mentored junior developers, providing guidance on best practices and design patterns
	Contributed to the development of the company's internal development standards and processes

	Senior Backend Developer, DEF Company (2016-2019)
	
	Developed and maintained a scalable backend system using Go and PostgreSQL
	Worked with a team of developers to implement new features and maintain the codebase
	Conducted code reviews, providing feedback and suggestions to improve code quality and maintainability
	Assisted in the development of the company's internal tools and processes

	Backend Developer, GHI Corporation (2014-2016)
	
	Developed and maintained a RESTful API using Python and Flask framework
	Worked with a team of developers to build and deploy a web-based analytics platform
	Contributed to the development of the company's internal coding standards and best practices

	SKILLS
	Programming languages: Go, Elixir, Python, Java
	Database management: PostgreSQL, MySQL
	Web development: RESTful APIs, microservices, Flask, Gin
	Cloud platforms: AWS, GCP
	Other: Git, Agile methodology, Scrum, JIRA
	
	EDUCATION
	Bachelor of Science in Computer Science, XYZ University (2014)
	
	CERTIFICATIONS
	
	Certified Kubernetes Administrator (CKA)
	Certified AWS Solutions Architect - Associate

	ACTIVITIES
	
	Volunteer, local tech community events
	Speaker, meetups and conferences on topics related to backend development and microservices`
)
