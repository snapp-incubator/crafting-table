package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const asciiArt = `
 ██████ ██████   █████  ███████ ████████ ██ ███    ██  ██████      ████████  █████  ██████  ██      ███████ 
██      ██   ██ ██   ██ ██         ██    ██ ████   ██ ██              ██    ██   ██ ██   ██ ██      ██      
██      ██████  ███████ █████      ██    ██ ██ ██  ██ ██   ███        ██    ███████ ██████  ██      █████   
██      ██   ██ ██   ██ ██         ██    ██ ██  ██ ██ ██    ██        ██    ██   ██ ██   ██ ██      ██      
 ██████ ██   ██ ██   ██ ██         ██    ██ ██   ████  ██████         ██    ██   ██ ██████  ███████ ███████ 
`

var rootCMD = &cobra.Command{
	Use:   "crafting-table",
	Short: "A repository for repository based struct",
}

// Execute executes the root command.
func Execute() {
	fmt.Print(asciiArt)
	if err := rootCMD.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	manifestCMD.AddCommand(applyCMD)
	rootCMD.AddCommand(manifestCMD)
}
