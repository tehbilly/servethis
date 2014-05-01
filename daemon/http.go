package daemon

import (
	"fmt"
	"net/http"
)

// Display index of enabled shares
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// Show list of contexts
	fmt.Fprint(w, "<table style='border: 0;'><tr><th>Context</th><th>Path</th></tr>")
	for _, v := range Daemon.FileServers {
		if v.Enabled {
			fmt.Fprintf(w, "<tr><td><a href='/%s/'>/%s/</a></td><td>%s</td></tr>", v.Context, v.Context, v.FilePath)
		}
	}
	fmt.Fprint(w, "</table>")
}

type ServeThisHTTPAdmin struct{}

// Display index PLUS administrative functions
func (s *ServeThisHTTPAdmin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// Show list of contexts
	fmt.Fprint(w, "<table style='border: 0;'><tr><th>Context</th><th>Path</th><th>Enable/Disable</th></tr>")
	for _, v := range Daemon.FileServers {
		if v.Enabled {
			fmt.Fprintf(w, "<tr><td><a href='/%s/'>/%s/</a></td><td>%s</td><td><a href='/_disable?ctx=%s'>Disable</a></tr>", v.Context, v.Context, v.FilePath, v.Context)
		} else {
			fmt.Fprintf(w, "<tr><td><a href='/%s/'>/%s/</a></td><td>%s</td><td><a href='/_enable?ctx=%s'>Enable</a></tr>", v.Context, v.Context, v.FilePath, v.Context)
		}
	}
	fmt.Fprint(w, "</table>")
	// Stop server
	// Add new path?
	// Usable stats?
}

type ServeHTTPHandler struct {
	Enabled  bool
	Context  string
	FilePath string
	handler  http.Handler
}

func (s *ServeHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Enabled {
		s.handler.ServeHTTP(w, r)
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func NewServeHTTPHandler(ctx string, path string) *ServeHTTPHandler {
	sh := &ServeHTTPHandler{
		Enabled:  true,
		Context:  ctx,
		FilePath: path,
		handler:  http.StripPrefix("/"+ctx+"/", http.FileServer(http.Dir(path))),
	}
	return sh
}
