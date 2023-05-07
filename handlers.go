package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	b64 "encoding/base64"

	"github.com/gorilla/mux"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type MethodsResponse struct {
	Status  string   `json:"status"`
	Methods []string `json:"methods"`
}

type ViewResponse struct {
	Status   string    `json:"status"`
	Torrents []Torrent `json:"torrents"`
}

type FilesResponse struct {
	Status string `json:"status"`
	Files  []File `json:"files"`
}

type PeersResponse struct {
	Status string `json:"status"`
	Peers  []Peer `json:"peers"`
}

type TrackersResponse struct {
	Status   string    `json:"status"`
	Trackers []Tracker `json:"trackers"`
}

func respond(p interface{}, statusCode int, w http.ResponseWriter) {
	bytes, err := json.Marshal(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func TemplateViewHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := []interface{}{"", "main", "d.hash=", "d.name=",
			"d.size_bytes=", "d.completed_bytes=", "d.up.rate=",
			"d.up.total=", "d.down.rate=", "d.down.total=",
			"d.message=", "d.is_active=", "d.is_open=",
			"d.state=", "d.state_changed=", "d.state_counter="}

		torrents, err := rt.DMulticall("main", args)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tpl := template.Must(template.ParseFiles("templates/torrents.html"))
		tpl.Execute(w, torrents)
	}
}

func HelloHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond(Response{
			Status:  "ok",
			Message: "🐢",
		}, http.StatusOK, w)
	}
}

func MethodsHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := rt.ListMethods()
		if err != nil {
			respond(Response{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusInternalServerError, w)
			return
		}
		respond(MethodsResponse{
			Status:  "ok",
			Methods: result,
		}, http.StatusOK, w)
	}
}

func LoadHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("error in load handler reading form: %s", err)
			respond(Response{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}
		defer file.Close()

		bytes := make([]byte, 0)
		_, err = file.Read(bytes)
		if err != nil {
			log.Printf("error in load handler loading bytes: %s", err)
			respond(Response{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}

		base64 := b64.StdEncoding.EncodeToString(bytes)
		err = rt.LoadRawStart(base64)
		if err != nil {
			log.Printf("error in load handler base64 encoding bytes: %s", err)
			respond(Response{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusInternalServerError, w)
			return
		}
		respond(Response{
			Status: "ok",
		}, http.StatusOK, w)
	}
}

func ViewHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		args := []interface{}{"", vars["view"],
			"d.hash=", "d.name=",
			"d.size_bytes=", "d.completed_bytes=", "d.up.rate=",
			"d.up.total=", "d.down.rate=", "d.down.total=",
			"d.message=", "d.is_active=", "d.is_open=",
			"d.state=", "d.state_changed=", "d.state_counter="}

		torrents, err := rt.DMulticall("main", args)
		if err != nil {
			log.Printf("error in view handler: %s", err)
			respond(Response{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}

		respond(ViewResponse{
			Status:   "ok",
			Torrents: torrents,
		}, http.StatusOK, w)
	}
}

func TorrentHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if vars["action"] == "stop" {
			err := rt.Stop(vars["hash"])
			if err != nil {
				log.Printf("error in action stop handler: %s", err)
				respond(Response{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusInternalServerError, w)
				return
			}
			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "start" {
			err := rt.Start(vars["hash"])
			if err != nil {
				log.Printf("error in action start handler: %s", err)
				respond(Response{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}
			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "files" {
			args := []interface{}{vars["hash"], "",
				"f.path=", "f.size_bytes=", "f.size_chunks=",
				"f.completed_chunks=", "f.frozen_path=", "f.priority=",
				"f.is_created=", "f.is_open="}

			files, err := rt.FMulticall(args)
			if err != nil {
				log.Printf("error in action files handler: %s", err)
				respond(Response{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}

			respond(FilesResponse{
				Status: "ok",
				Files:  files,
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "peers" {
			args := []interface{}{vars["hash"], "",
				"p.id=", "p.address=", "p.port=",
				"p.banned=", "p.client_version=", "p.completed_percent=",
				"p.is_encrypted=", "p.is_incoming=", "p.is_obfuscated=",
				"p.peer_rate=", "p.peer_total=", "p.up_rate=", "p.up_total="}

			peers, err := rt.PMulticall(args)
			if err != nil {
				log.Printf("error in action peers handler: %s", err)
				respond(Response{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}

			respond(PeersResponse{
				Status: "ok",
				Peers:  peers,
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "trackers" {
			args := []interface{}{vars["hash"], "",
				"t.id=", "t.type=", "t.url=",
				"t.activity_time_last=", "t.activity_time_next=", "t.can_scrape=",
				"t.is_usable=", "t.is_enabled=", "t.failed_counter=",
				"t.failed_time_last=", "t.failed_time_next=", "t.is_busy=",
				"t.is_open=",
			}

			trackers, err := rt.TMulticall(args)
			if err != nil {
				log.Printf("error in action trackers handler: %s", err)
				respond(Response{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}

			respond(TrackersResponse{
				Status:   "ok",
				Trackers: trackers,
			}, http.StatusOK, w)
			return
		}

	}
}
