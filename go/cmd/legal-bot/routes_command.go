package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var routesJSON bool

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Show the active Legal-Bot route catalog and legacy surface inventory",
	RunE:  runRoutes,
}

func init() {
	routesCmd.Flags().BoolVar(&routesJSON, "json", false, "Print the route catalog as JSON")
	rootCmd.AddCommand(routesCmd)
}

func runRoutes(cmd *cobra.Command, args []string) error {
	if routesJSON {
		data, err := marshalRouteSnapshotJSON()
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	}

	snapshot := buildRouteSnapshot()
	w := cmd.OutOrStdout()
	fmt.Fprintln(w, "Active commands:")
	for _, name := range snapshot.ActiveCommands {
		fmt.Fprintf(w, "- %s\n", name)
	}

	fmt.Fprintln(w, "\nReview modes:")
	for _, mode := range snapshot.ReviewModes {
		fmt.Fprintf(w, "- %s -> %s\n", mode.Name, mode.PipelineKey)
	}

	fmt.Fprintln(w, "\nReview document types:")
	for _, document := range snapshot.ReviewDocumentTypes {
		fmt.Fprintf(w, "- %s -> %s\n", document.Name, document.MapsTo)
	}

	fmt.Fprintln(w, "\nWorkflow types:")
	for _, workflow := range snapshot.WorkflowTypes {
		fmt.Fprintf(w, "- %s -> %s\n", workflow.Type, workflow.PipelineKey)
	}

	fmt.Fprintln(w, "\nConditional specialists:")
	for _, rule := range snapshot.ConditionalRules {
		fmt.Fprintf(w, "- %s (%s)\n", rule.Persona, rule.Trigger)
	}

	fmt.Fprintln(w, "\nLegacy surfaces:")
	for _, surface := range snapshot.LegacySurfaces {
		fmt.Fprintf(w, "- %s [%s] -> %s\n", surface.Path, surface.Status, surface.Replacement)
	}

	return nil
}
