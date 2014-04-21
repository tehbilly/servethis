package daemon

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	foreground bool
	port       int
)

var DaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Runs the servethis daemon. This generally does not need to be called directly.",
	Long: `The daemon is what actually serves the directory contents. The base servethis command (and related ones)
will run this command if it is required.

Please do not run this command directly unless you know what you are doing!`,
	Run: func(cmd *cobra.Command, args []string) {
		dKey := "_SERVETHIS_DAEMON"
		dVal := "true"
		// Spawn a child process if we are not already daemonized
		if !foreground && os.Getenv(dKey) != dVal {
			path, err := filepath.Abs(os.Args[0])
			if err != nil {
				log.Fatal(err)
			}

			cmd := exec.Command(path, os.Args[1:]...)
			cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", dKey, dVal))
			if err = cmd.Start(); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}

		// Start RPC server

		// Start web server
		http.HandleFunc("/", statusHandler)
		http.ListenAndServe(":8000", nil)
	},
}

func init() {
	DaemonCmd.Flags().BoolVar(&foreground, "foreground", false, "Run the daemon in the foreground.")
	DaemonCmd.Flags().IntVar(&port, "port", 8000, "The port the daemon should serve content on")
}
