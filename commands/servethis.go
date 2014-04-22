package commands

import (
	"bytes"
	"errors"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tehbilly/servethis/daemon"
)

var ServeThisCmd = &cobra.Command{
	Use:   "servethis",
	Short: "Serve a directory via http!",
	Long:  "servethis is a utility for serving files via http.",
	Run: func(cmd *cobra.Command, args []string) {
		for i := 1; i <= 5; i++ {
			// Dial the daemon
			conn, err := rpc.Dial("tcp", "127.0.0.1:8034")
			if err != nil {
				derr := StartDaemon()
				if derr != nil {
					fmt.Printf("Error attempting to start daemon during attempt %d. Error message: %v\n", i, derr)
				}
				continue
			} else {
				var url string
				rpcArgs := &daemon.AddArgs{Context: ServeContext, Path: ServePath}
				err = conn.Call("ServeThis.AddHandler", *rpcArgs, &url)
				if err != nil {
					fmt.Println("Error adding path to daemon:", err)
					return
				}

				fmt.Printf("%s is being served up at: %s\n", ServePath, url)
				break
			}
		}
	},
}

var (
	ServePath    string
	ServeContext string
	servePort    int
)

// Stuff we gonna need pretty much urrywhurr
func init() {
	// Get the current working directory for defaults
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	var ctx string
	abs, err := filepath.Abs(cwd)
	if err != nil {
		abs = "."
		ctx = "serve"
	} else {
		ctx = filepath.Base(abs)
	}

	// Set up global (persisitent) variables
	ServeThisCmd.Flags().StringVar(&ServePath, "path", cwd, "Specify the path to serve. Defaults to current directory.")
	ServeThisCmd.Flags().StringVarP(&ServeContext, "context", "c", ctx, "Context root for the served directory. Defaults to directory name (if available).")
	ServeThisCmd.Flags().IntVar(&servePort, "port", 8000, "Port the daemon will listen on (if not already started).")

	// Add subcommands
	ServeThisCmd.AddCommand(versionCmd)
	ServeThisCmd.AddCommand(daemon.DaemonCmd)
	ServeThisCmd.AddCommand(StopCmd)
	ServeThisCmd.AddCommand(RemoveCmd)
}

func StartDaemon() error {
	// Start the daemon if it ain't running
	_, err := rpc.Dial("tcp", "127.0.0.1:8034")
	if err != nil { // Can't reach it? Well... start it!
		fmt.Println("Daemon is not currently running. Starting daemon...")
		cmd := exec.Command(os.Args[0], "daemon")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Println("Unable to start daemon!", err)
			return errors.New("Daemon did not start")
		}
	}
	return nil
}
