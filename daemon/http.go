package daemon

import (
	"fmt"
	"net/http"
)

// Display status of the servethis daemon/server, and exposed paths as a browseable index
func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Serve something simple for now
	fmt.Fprintf(w, "This is a simple status handler! There will one day be an index below. Yup.\nMessage: %s", "Herro")
}
