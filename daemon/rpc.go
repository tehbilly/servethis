package daemon

import "errors"

type RpcArgs struct {
	Context string
	Path    string
}

// ServeThis is the container for the RPC methods
type ServeThis int

// AddPath will add a new path to the HTTP handler
// A string with the full accessible URL will be returned on success
func (s *ServeThis) AddPath(args *RpcArgs, reply *string) error {
	// Make sure we actually have valid arguments
	if args.Context == "" || args.Path == "" {
		return errors.New("Context and Path must not contain empty strings.")
	}
	// Possibly handle a relative path here?

	// daemon.AddHandler(Context, AbsolutePath)
	// return full handler URL
}
