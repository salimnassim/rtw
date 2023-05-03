package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	transport := &http.Transport{}

	if os.Getenv("BASIC_USERNAME") != "" && os.Getenv("BASIC_PASSWORD") != "" {
		transport.RegisterProtocol("https",
			newBasicAuthTransport(
				os.Getenv("BASIC_USERNAME"),
				os.Getenv("BASIC_PASSWORD"),
			),
		)
	}

	rtorrent, err := NewRtorrent(RtorrentConfig{
		URL:       os.Getenv("URL"),
		Transport: transport,
	})
	if err != nil {
		log.Fatalf("unable to create rtorrent instance: %v", err)
	}

	defer rtorrent.client.Close()

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/view/{view}", ViewHandler(rtorrent))
	s.HandleFunc("/torrent/{hash}/{action}", TorrentHandler(rtorrent))

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              os.Getenv("BIND_ADDRESS"),
		Handler:           r,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("unable start http server: %s", err)
	}

}

type basicAuthTransport struct {
	Username string
	Password string
}

func newBasicAuthTransport(username, password string) *basicAuthTransport {
	return &basicAuthTransport{
		Username: username,
		Password: password,
	}
}

func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return http.DefaultTransport.RoundTrip(req)
}
