package cmd

import (
	"github.com/open-identity/roles-manager/cmd/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server and serves the HTTP REST API",
	Long:  "This command opens a network port and listens to HTTP/2 API requests.",
	Run:   server.RunServe(logger, Version, Commit, Date),
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
