package daemon

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	Daemon     *serveDaemon
	hostname   string
	foreground bool
	httpport   int
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

		// Start our shutdown checker!
		go shutdownOnEmptyCheck()

		// Start RPC server
		go startRPC()

		// Our status/index handler
		http.HandleFunc("/", statusHandler)

		// Shouldn't hurt if stdout is piped to nowhere
		fmt.Printf("Daemon started successfully. Starting HTTP server at http://%s:%d\n", hostname, httpport)

		// Start web server
		go http.ListenAndServe(":"+strconv.Itoa(httpport), nil)

		// Handle stop requests or OS signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, os.Kill)
		select {
		case <-sigChan:
			fmt.Println("Caught signal! Exiting...")
			os.Exit(0)
		case <-Daemon.QuitChan:
			fmt.Println("Got message on quit channel. Exiting.")
			os.Exit(0)
		}
	},
}

func init() {
	// Save hostname for building URLs
	host, _ := os.Hostname()
	hostname = strings.ToLower(host)

	// Establish our Daemon struct
	Daemon = NewServeDaemon()

	// Flags specific to this cobra.Command
	DaemonCmd.Flags().BoolVar(&foreground, "foreground", false, "Run the daemon in the foreground.")
	DaemonCmd.Flags().IntVar(&httpport, "httpport", 8000, "Port the daemon will listen on")
}

// A simple function that will check to make sure we exist for a reason. If we have 0 contexts, shut down
func shutdownOnEmptyCheck() {
	for {
		<-time.After(30 * time.Second)
		if len(Daemon.FileServers) == 0 {
			os.Exit(1)
		}
	}
}

// Register our RPC service and start a listener/handler for it
func startRPC() {
	rpc.RegisterName("ServeThis", new(ServeThisRPC))
	l, err := net.Listen("tcp", "127.0.0.1:8034")
	defer l.Close()

	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			continue // TODO Log this somewhere?
		}
		go rpc.ServeConn(conn)
	}
}

func NewServeDaemon() *serveDaemon {
	sd := &serveDaemon{
		FileServers: make(map[string]string),
		QuitChan:    make(chan bool, 1),
	}
	return sd
}

type serveDaemon struct {
	FileServers map[string]string // map[context]path
	QuitChan    chan bool
}

func (s *serveDaemon) Shutdown() {
	s.QuitChan <- true
}

func AddFileServer(context, path string) (string, error) {
	var ctx string // This is the actual context this handler is registered at
	ctx = strings.TrimLeft(context, "/")
	ctx = strings.TrimRight(ctx, "/")

	// Do we already have this context being used? Make a unique name instead of erroring.
	if _, contains := Daemon.FileServers[ctx]; contains {
		ctxIndex := 1
		for {
			tmpCtx := fmt.Sprintf("%s-%d", ctx, ctxIndex)
			if val, contains := Daemon.FileServers[tmpCtx]; contains && val != "" {
				ctxIndex++
				continue
			} else {
				ctx = tmpCtx
				break
			}
		}
	}

	// is the context already used? Eh, for now we'll just write over the old one
	Daemon.FileServers[ctx] = path // We want this as just a name in here
	ctx = "/" + ctx + "/"
	http.Handle(ctx, http.StripPrefix(ctx, http.FileServer(http.Dir(path))))

	serveURL := "http://" + hostname + ":" + strconv.Itoa(httpport) + ctx
	return serveURL, nil
}

func RemoveFileServer(context string) {
	http.HandleFunc("/"+context+"/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "There is no longer any content being served at: %s", context)
	})
	Daemon.FileServers[context] = ""
}
