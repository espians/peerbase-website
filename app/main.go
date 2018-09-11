// Public Domain (-) 2018-present, The Peerbase Website Authors.
// See the Peerbase Website UNLICENSE file for details.

package main

import (
	"net/http"
	"strings"

	"google.golang.org/appengine"
)

// TODO(tav): Switch over to using `mod` and Go module proxy.
var goMeta = []byte(`<!doctype html>
<meta name="go-import" content="peerbase.net/go git https://github.com/peerbase/peerbase">
<meta name="go-source" content="peerbase.net/go https://github.com/peerbase/peerbase https://github.com/peerbase/peerbase/tree/master{/dir} https://github.com/peerbase/peerbase/blob/master{/dir}/{file}#L{line}">`)

func handle(w http.ResponseWriter, r *http.Request) {
	if !appengine.IsDevAppServer() {
		if r.Host != "peerbase.net" || r.URL.Scheme != "https" {
			url := r.URL
			url.Host = "peerbase.net"
			url.Scheme = "https"
			http.Redirect(w, r, url.String(), 301)
			return
		}
	}
	split := strings.Split(r.URL.Path, "/")
	if len(split) >= 2 && split[1] == "go" {
		handleGo(w, r, split[2:])
		return
	}
	http.Redirect(w, r, "https://github.com/peerbase/peerbase", 302)
}

func handleGo(w http.ResponseWriter, r *http.Request, path []string) {
	if r.URL.Query().Get("go-get") != "" {
		w.Write(goMeta)
		return
	}
	if r.URL.RawQuery != "" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "https://godoc.org/peerbase.net/go/"+strings.Join(path, "/"), 302)
}

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
}
