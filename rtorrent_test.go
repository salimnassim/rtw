package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMulticallSystem(t *testing.T) {
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
		t.Error(err)
	}

	args := []interface{}{
		[]interface{}{
			SystemCall{
				MethodName: "throttle.global_down.total",
				Params:     []string{""},
			},
			SystemCall{
				MethodName: "throttle.global_up.total",
				Params:     []string{""},
			},
		},
	}

	result, err := rtorrent.SystemMulticall(args)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%v", result)
}
