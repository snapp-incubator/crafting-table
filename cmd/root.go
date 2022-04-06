package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "crafting-table",
	Short: "A generator for repository based struct",
}

// Execute executes the root command.
func Execute() {
	printAsciiArt()
	if err := rootCMD.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCMD.AddCommand(generateCMD)
}
