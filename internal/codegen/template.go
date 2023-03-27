package codegen

import (
	_ "embed"
)

//go:embed templates/template.tmpl
var templateStr string
