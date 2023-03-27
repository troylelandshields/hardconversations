package main

import (
	"os"

	"github.com/troylelandshields/hardconversations/internal/cmd"
)

func main() {
	os.Exit(cmd.Do(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
