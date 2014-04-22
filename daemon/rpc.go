package daemon

import (
	"errors"
	"strings"
)

type AddArgs struct {
	Context string
	Path    string
}

// ServeThis is the container for the RPC methods
type ServeThisRPC int

// AddHandler will add a new path to the HTTP handler
// A string with the full accessible URL will be returned on success
func (s *ServeThisRPC) AddHandler(args *AddArgs, reply *string) error {
	// Make sure we actually have valid arguments
	if args.Context == "" || args.Path == "" {
		return errors.New("Context and Path must not contain empty strings.")
	}
	// Possibly handle a relative path here? Gotta be safe and all.

	// daemon.AddHandler(Context, AbsolutePath)
	// return full handler URL
	url, err := AddFileServer(args.Context, args.Path)
	if err != nil {
		return err
	}

	*reply = url
	return nil
}

func (s *ServeThisRPC) RemoveHandler(context *string, reply *string) error {
	ctx := strings.TrimLeft(*context, "/")
	ctx = strings.TrimRight(ctx, "/")

	if val, contains := Daemon.FileServers[ctx]; contains && val != "" {
		RemoveFileServer(ctx)
		*reply = "No longer serving content for: " + ctx
	} else {
		*reply = "Unable to remove handler. Nothing is currently served at: " + ctx
	}

	return nil
}

func (s *ServeThisRPC) StopDaemon(_, _ *struct{}) error {
	Daemon.Shutdown()
	return nil
}
