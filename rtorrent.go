package main

import (
	"net/http"
	"reflect"

	"github.com/kolo/xmlrpc"
)

type Torrent struct {
	Hash           string `xmlrpc:"d.hash=" json:"hash"`
	Name           string `xmlrpc:"d.name=" json:"name"`
	Size           int64  `xmlrpc:"d.size_bytes=" json:"size_bytes"`
	CompletedBytes int64  `xmlrpc:"d.completed_bytes=" json:"completed_bytes"`
	UploadRate     int64  `xmlrpc:"d.up.rate=" json:"upload_rate"`
	UploadTotal    int64  `xmlrpc:"d.up.total=" json:"upload_total"`
	DownloadRate   int64  `xmlrpc:"d.down.rate=" json:"download_rate"`
	DownloadTotal  int64  `xmlrpc:"d.down.total=" json:"download_total"`
	Message        string `xmlrpc:"d.message=" json:"message"`
	IsActive       int64  `xmlrpc:"d.is_active=" json:"is_active"`
	IsOpen         int64  `xmlrpc:"d.is_open=" json:"is_open"`
	IsHashing      int64  `xmlrpc:"d.is_hash_checking=" json:"is_hashing"`
	Leechers       int64  `xmlrpc:"d.peers_accounted=" json:"leechers"`
	Seeders        int64  `xmlrpc:"d.peers_complete=" json:"seeders"`
	Custom1        string `xmlrpc:"d.custom1=" json:"custom1"`
	Custom2        string `xmlrpc:"d.custom2=" json:"custom2"`
	Custom3        string `xmlrpc:"d.custom3=" json:"custom3"`
	Custom4        string `xmlrpc:"d.custom4=" json:"custom4"`
	Custom5        string `xmlrpc:"d.custom5=" json:"custom5"`
}

type RtorrentConfig struct {
	URL       string
	Transport http.RoundTripper
}

type Rtorrent struct {
	client *xmlrpc.Client
}

func NewRtorrent(config RtorrentConfig) (*Rtorrent, error) {
	xmlrpcClient, err := xmlrpc.NewClient(config.URL, config.Transport)
	if err != nil {
		return nil, err
	}

	rtorrent := &Rtorrent{
		client: xmlrpcClient,
	}
	return rtorrent, nil
}

func (rt *Rtorrent) ListMethods() (interface{}, error) {
	var result interface{}
	err := rt.client.Call("system.listMethods", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rt *Rtorrent) DMulticall(view string, args interface{}) ([]Torrent, error) {
	var result interface{}
	err := rt.client.Call("d.multicall2", args, &result)
	if err != nil {
		return nil, err
	}

	torrents := make([]Torrent, 0)
	for _, outer := range result.([]interface{}) {
		torrent := &Torrent{}
		for idx := 2; idx < len(args.([]interface{})); idx++ {
			ref := outer.([]interface{})[idx-2]
			fname := args.([]interface{})[idx].(string)
			vo := reflect.ValueOf(torrent)
			el := vo.Elem()
			for i := 0; i < el.NumField(); i++ {
				field := el.Type().Field(i)
				if fname == field.Tag.Get("xmlrpc") {
					if ref == nil {
						continue
					}
					el.Field(i).Set(reflect.ValueOf(ref))
				}
			}

		}
		torrents = append(torrents, *torrent)
	}

	return torrents, nil
}
