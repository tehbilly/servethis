package daemon

import (
	"fmt"
	"net/http"
)

// Display status of the servethis daemon/server, and exposed paths as a browseable index
func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// Show list of contexts
	fmt.Fprint(w, "<table style='border: 0;'><tr><th>Context</th><th>Path</th></tr>")
	for k, v := range Daemon.FileServers {
		fmt.Fprintf(w, "<tr><td><a href='/%s/'>/%s/</a></td><td>%s</td></tr>", k, k, v)
	}
	fmt.Fprint(w, "</table>")
}
