// Package main implements a CLI tool that scans Go source files for
// calls to the i18n.Bundle.Translate method and extracts the format
// strings to generate JSON translation files. It automatically detects
// bundle configurations from the source code.
package main

import (
	"fmt"
	"os"

	"github.com/epkgs/i18n/i18ntool/cmd"
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

	rootCmd.AddCommand(cmd.ExtractCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
