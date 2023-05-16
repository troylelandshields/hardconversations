package codegen

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/troylelandshields/hardconversations/internal/config"
)

func ExecuteTemplate(convo config.Conversation) ([]byte, error) {
	tmpl := template.Must(template.New("tmpl").Parse(templateStr))

	var buf bytes.Buffer
	tmplData := buildTmplCtx(convo)
	err := tmpl.Execute(&buf, tmplData)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func buildTmplCtx(convo config.Conversation) tmplCtx {
	imports := make(map[string]struct{})

	for _, q := range convo.Questions {
		if q.InputParsed.ImportPath != "" {
			imports[q.InputParsed.ImportPath] = struct{}{}
		}
		if q.OutputParsed.ImportPath != "" {
			imports[q.OutputParsed.ImportPath] = struct{}{}
		}
	}

	var importList []string
	for k := range imports {
		importList = append(importList, k)
	}

	return tmplCtx{
		Conversation: convo,
		Package:      strings.Trim(convo.Path, "./"),
		Imports:      importList,
	}
}

type tmplCtx struct {
	config.Conversation
	Package string
	Imports []string
}
