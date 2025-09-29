// Package main implements a CLI tool that scans Go source files for
// calls to the i18n.Bundle.Translate method and extracts the format
// strings to generate JSON translation files. It automatically detects
// bundle configurations from the source code.
package main

import (
	"fmt"
	"os"

	"github.com/epkgs/i18n/i18ntool/internal"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "i18n",
		Short: "i18n is a CLI tool for managing translation files",
		Long: `i18n is a CLI tool that scans Go source files for calls to the 
i18n.Bundle.Translate method and extracts the format strings to generate 
JSON translation files. It automatically detects bundle configurations 
from the source code.`,
	}

	rootCmd.AddCommand(extractCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func extractCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract translation definitions by scanning source code",
		RunE: func(cmd *cobra.Command, args []string) error {

			searchPath, err := cmd.Flags().GetString("path")
			if err != nil {
				return err
			}

			langs, _ := cmd.Flags().GetStringSlice("lang")

			g := internal.NewGenerator(searchPath)

			if err := g.CollectBundles(); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			if err := g.GenerateTranslations(langs...); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Translation files generated successfully")

			return nil
		},
	}

	cmd.Flags().StringP("path", "p", ".", "Path to search for Go source files")
	cmd.Flags().StringSliceP("lang", "l", []string{}, "Languages to generate translations for")

	return cmd
}
