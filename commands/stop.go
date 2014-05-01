package commands

import (
	"fmt"
	"net/rpc"

	"github.com/spf13/cobra"
)

// StopCommand will send a stop/shutdown request to the daemon
var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop serving entirely. Shuts down the daemon if running",
	Run: func(cmd *cobra.Command, args []string) {
		// Try to dial the daemon
		conn, err := rpc.Dial("tcp", "127.0.0.1:8034")
		if err != nil {
			fmt.Println("Daemon is not reachable. Can't shut down something that ain't running!")
			return
		}

		var rargs struct{}
		var resp struct{}
		err = conn.Call("ServeThis.StopDaemon", rargs, &resp)
		if err != nil {
			fmt.Println("Response from daemon:", err)
		} else {
			fmt.Println("Daemon should now be stopped!")
		}
	},
}
