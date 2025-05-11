// cmd/ghostpost/images.go

package main

import "github.com/spf13/cobra"

func imagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "images",
		Short: "Image operations (placeholder)",
		Run: func(_ *cobra.Command, _ []string) {
			println("Not implemented yet—coming soon.")
		},
	}
}
