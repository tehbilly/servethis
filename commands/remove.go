package commands

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:   "remove [context]",
	Short: "Remove a handler identified by [context].",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Invalid number of arguments. Specify the context to remove. Usage:", os.Args[0], "context")
			return
		}

		// Try to dial the daemon
		conn, err := rpc.Dial("tcp", "127.0.0.1:8034")
		if err != nil {
			fmt.Println("Daemon is not reachable. Are you sure it's running?\n", err)
		}

		context := args[0]
		var response string
		err = conn.Call("ServeThis.RemoveHandler", context, &response)
		if err != nil {
			fmt.Println("Daemon returned an error:", err)
		} else {
			fmt.Println(response)
		}
	},
}
