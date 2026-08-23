/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/eneiss/gar/internal/cli/extract"
)

func init() {
	rootCmd.AddCommand(extract.NewCommand())
	// extractCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
