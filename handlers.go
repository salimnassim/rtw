package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ViewResponse struct {
	Status   string    `json:"status"`
	Torrents []Torrent `json:"torrents,omitempty"`
}

func response(p interface{}, statusCode int, w http.ResponseWriter) {
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
			response(ErrorResponse{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}

		response(ViewResponse{
			Status:   "ok",
			Torrents: torrents,
		}, http.StatusOK, w)
	}
}
