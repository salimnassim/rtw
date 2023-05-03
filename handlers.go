package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type ViewResponse struct {
	Status   string    `json:"status"`
	Torrents []Torrent `json:"torrents,omitempty"`
}

type TorrentResponse struct {
	Status string `json:"status"`
	Files  []File `json:"files,omitempty"`
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

func ViewHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		args := []interface{}{"", vars["view"], "d.hash=", "d.name=",
			"d.size_bytes=", "d.completed_bytes=", "d.up.rate=",
			"d.up.total=", "d.down.rate=", "d.down.total=",
			"d.message=", "d.is_active=", "d.is_open="}

		torrents, err := rt.DMulticall("main", args)
		if err != nil {
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
				respond(Response{
					Status:  "error",
					Message: fmt.Sprintf("unable to stop torrent: %v", err),
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
				respond(Response{
					Status:  "error",
					Message: fmt.Sprintf("unable to start torrent: %v", err),
				}, http.StatusBadRequest, w)
				return
			}
			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "files" {
			args := []interface{}{vars["hash"], "", "f.path=", "f.size_chunks=",
				"f.completed_chunks=", "f.frozen_path=", "f.priority=",
				"f.is_created=", "f.is_open="}

			files, err := rt.FMulticall(args)
			if err != nil {
				respond(Response{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}

			respond(TorrentResponse{
				Status: "ok",
				Files:  files,
			}, http.StatusOK, w)
			return
		}

	}
}
