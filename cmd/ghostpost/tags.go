// cmd/ghostpost/tags.go

package main

import "github.com/spf13/cobra"

func tagsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tags",
		Short: "Tag operations (placeholder)",
		Run: func(_ *cobra.Command, _ []string) {
			println("Not implemented yet—coming soon.")
		},
	}
}
