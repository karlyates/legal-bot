package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:     "legal-bot",
	Short:   "Legal-Bot - local-first legal analysis, workflow, and drafting support",
	Long:    "Legal-Bot builds matter knowledge from local source files and runs practical legal analysis, triage, drafting, authority-check, and draft-review workflows.",
	Version: version,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
