package compiler

import (
	"context"
	"path/filepath"

	"github.com/troylelandshields/hardconversations/internal/codegen"
	"github.com/troylelandshields/hardconversations/internal/config"
)

type Compiler struct {
	conf config.Conversation
}

func NewCompiler(conf config.Conversation) *Compiler {
	c := &Compiler{conf: conf}
	return c
}

func (c *Compiler) Compile(ctx context.Context) ([]CompileResult, error) {
	var results []CompileResult

	b, err := codegen.ExecuteTemplate(c.conf)
	if err != nil {
		return nil, err
	}

	results = append(results, CompileResult{
		Name:     filepath.Join(c.conf.Path, "client.gen.go"),
		Contents: b,
	})

	return results, nil
}

type CompileResult struct {
	Name     string
	Contents []byte
}

type tmplCtx struct {
}
