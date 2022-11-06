package cmd

import (
	"github.com/snapp-incubator/crafting-table/internal/server"
	"github.com/spf13/cobra"
)

var port string

var serveCMD = &cobra.Command{
	Use:   "serve",
	Short: "Start generating repository",
	Run:   serve,
}

func init() {
	serveCMD.Flags().StringVarP(&port, "port", "p", "7628", "port of server")
}

func serve(_ *cobra.Command, _ []string) {
	if port == "" {
		panic("port is not set")
	}
	server.Serve(port)
}
