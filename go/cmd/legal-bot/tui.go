package main

import (
	"github.com/jdonohoo/legal-bot/go/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive Bubble Tea TUI utility surface",
	Long:  "Start the Legal-Bot terminal UI for legacy-oriented discovery, council, historian, and settings utilities. This is not the primary family-law workflow front door.",
	RunE:  runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(cmd *cobra.Command, args []string) error {
	agentsDir := resolveAgentsDir()

	// Find project root
	projectRoot := ""
	if agentsDir != "agents" && len(agentsDir) > len("/agents") {
		projectRoot = agentsDir[:len(agentsDir)-len("/agents")]
	}

	return tui.Run(projectRoot, agentsDir, version)
}
