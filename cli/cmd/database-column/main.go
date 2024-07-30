package main

import (
	"log"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "column",
		Short: "Generate column definition files from GORM models",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatalf("Usage: column <model_path>")
			}
			modelPath := args[0]
			err := generateColumnFiles(modelPath)
			if err != nil {
				log.Fatalf("Error generating column files: %v", err)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
