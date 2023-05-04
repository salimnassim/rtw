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
		return
	}

	defer rtorrent.client.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", TemplateViewHandler(rtorrent))

	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/load", LoadHandler(rtorrent)).Methods("POST")
	s.HandleFunc("/methods", MethodsHandler(rtorrent))
	s.HandleFunc("/view/{view}", ViewHandler(rtorrent))
	s.HandleFunc("/torrent/{hash}/{action}", TorrentHandler(rtorrent))
	s.Use(CorsMiddleware)

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              os.Getenv("BIND_ADDRESS"),
		Handler:           r,
	}

	log.Printf("listen address: http://%s", os.Getenv("BIND_ADDRESS"))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("server failure: %s", err)
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
