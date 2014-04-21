package commands

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tehbilly/servethis/daemon"
)

var ServeThisCmd = &cobra.Command{
	Use:   "servethis",
	Short: "Serve a directory via http!",
	Long:  "servethis is a utility for serving files via http.",
	Run: func(cmd *cobra.Command, args []string) {
		// Add a basic do-nothing handler
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Println("Hey")
		})

		// We want the hostname to make a nice easy to copy link
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "0.0.0.0"
		}
		hostname = strings.ToLower(hostname)

		ctx := "/" + ServeContext + "/"

		log.Printf("Starting http server to serve '%s' at:\nhttp://%s:%s%s", ServePath, hostname, "8000", ctx)
		//fileHandler := http.FileServer(http.Dir(ServePath))
		//wrappedHandler := AccessLoggingHandler(fileHandler)
		http.Handle(ctx, http.StripPrefix(ctx, http.FileServer(http.Dir("."))))
		log.Fatal(http.ListenAndServe(":8000", nil))
	},
}

var (
	ServePath     string
	ServeContext  string
	ExposeOnIndex bool
	// The following are on the base command because this will call the server command with appropriate flags
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
	ServeThisCmd.PersistentFlags().StringVar(&ServePath, "path", cwd, "Specify the path to serve. Defaults to current directory.")
	ServeThisCmd.PersistentFlags().StringVarP(&ServeContext, "context", "c", ctx, "Context root for the served directory. Defaults to directory name (if available).")
	ServeThisCmd.PersistentFlags().BoolVarP(&ExposeOnIndex, "index", "i", true, "Show this listing on the index?") // Should this default to false?

	// Add subcommands
	ServeThisCmd.AddCommand(versionCmd)
	ServeThisCmd.AddCommand(daemon.DaemonCmd)
}

// A bit excessive? Perhaps. But I enjoy seeing what's going on while the server is running.
func AccessLoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RemoteAddr, "=>", r.RequestURI)
		h.ServeHTTP(w, r)
	})
}
