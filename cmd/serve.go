package cmd

import (
	"github.com/snapp-incubator/crafting-table/internal/server"
	"github.com/spf13/cobra"
)

var serveCMD = &cobra.Command{
	Use:   "serve",
	Short: "Start generating repository",
	Run:   serve,
}

func serve(_ *cobra.Command, _ []string) {
	server.Serve()
}
