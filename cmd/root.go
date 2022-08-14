package cmd

import (
	"fmt"
	"os"

	"github.com/snapp-incubator/crafting-table/internal/assets"

	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "crafting-table",
	Short: "A repository for repository based struct",
}

// Execute executes the root command.
func Execute() {
	assets.PrintAsciiArt()
	if err := rootCMD.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCMD.AddCommand(generateCMD, uiCMD)
}
