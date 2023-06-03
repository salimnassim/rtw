package main

import (
	"encoding/base64"
	"net/http"
	"reflect"

	"github.com/kolo/xmlrpc"
)

type Torrent struct {
	Hash           string `xmlrpc:"d.hash=" json:"hash"`
	Name           string `xmlrpc:"d.name=" json:"name"`
	SizeBytes      int64  `xmlrpc:"d.size_bytes=" json:"size_bytes"`
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
	State          int64  `xmlrpc:"d.state=" json:"state"`
	StateChanged   int64  `xmlrpc:"d.state_changed=" json:"state_changed"`
	StateCounter   int64  `xmlrpc:"d.state_counter=" json:"state_counter"`
	Priority       int64  `xmlrpc:"d.priority=" json:"priority"`
	Custom1        string `xmlrpc:"d.custom1=" json:"custom1"`
	Custom2        string `xmlrpc:"d.custom2=" json:"custom2"`
	Custom3        string `xmlrpc:"d.custom3=" json:"custom3"`
	Custom4        string `xmlrpc:"d.custom4=" json:"custom4"`
	Custom5        string `xmlrpc:"d.custom5=" json:"custom5"`
}

type File struct {
	Path            string `xmlrpc:"f.path=" json:"path"`
	Size            int64  `xmlrpc:"f.size_bytes=" json:"size"`
	SizeChunks      int64  `xmlrpc:"f.size_chunks=" json:"size_chunks"`
	CompletedChunks int64  `xmlrpc:"f.completed_chunks=" json:"completed_chunks"`
	FrozenPath      string `xmlrpc:"f.frozen_path=" json:"frozen_path"`
	Priority        int64  `xmlrpc:"f.priority=" json:"priority"`
	IsCreated       int64  `xmlrpc:"f.is_created=" json:"is_created"`
	IsOpen          int64  `xmlrpc:"f.is_open=" json:"is_open"`
}

type Peer struct {
	PeerID           string `xmlrpc:"p.id=" json:"peer_id"`
	Address          string `xmlrpc:"p.address=" json:"address"`
	Port             int64  `xmlrpc:"p.port=" json:"port"`
	Banned           int64  `xmlrpc:"p.banned=" json:"banned"`
	ClientVersion    string `xmlrpc:"p.client_version=" json:"client_version"`
	CompletedPercent int64  `xmlrpc:"p.completed_percent=" json:"completed_percent"`
	IsEncrypted      int64  `xmlrpc:"p.is_encrypted=" json:"is_encrypted"`
	IsIncoming       int64  `xmlrpc:"p.is_incoming=" json:"is_incoming"`
	IsObfuscated     int64  `xmlrpc:"p.is_obfuscated=" json:"is_obfuscated"`
	PeerRate         int64  `xmlrpc:"p.peer_rate=" json:"peer_rate"`
	PeerTotal        int64  `xmlrpc:"p.peer_total=" json:"peer_total"`
	UploadRate       int64  `xmlrpc:"p.up_rate=" json:"up_rate"`
	UploadTotal      int64  `xmlrpc:"p.up_total=" json:"up_total"`
}

type Tracker struct {
	TrackerID        string `xmlrpc:"t.id=" json:"tracker_id"`
	ActivityTimeLast int64  `xmlrpc:"t.activity_time_last=" json:"activity_time_last"`
	ActivityTimeNext int64  `xmlrpc:"t.activity_time_next=" json:"activity_time_next"`
	CanScrape        int64  `xmlrpc:"t.can_scrape=" json:"can_scrape"`
	IsUsable         int64  `xmlrpc:"t.is_usable=" json:"t.is_usable"`
	IsEnabled        int64  `xmlrpc:"t.is_enabled=" json:"is_enabled"`
	FailedCounter    int64  `xmlrpc:"t.failed_counter=" json:"failed_counter"`
	FailedTimeLast   int64  `xmlrpc:"t.failed_time_last=" json:"failed_time_last"`
	FailedTimeNext   int64  `xmlrpc:"t.failed_time_next=" json:"failed_time_next"`
	IsBusy           int64  `xmlrpc:"t.is_busy=" json:"is_busy"`
	IsOpen           int64  `xmlrpc:"t.is_open=" json:"is_open"`
	Type             int64  `xmlrpc:"t.type=" json:"type"`
	URL              string `xmlrpc:"t.url=" json:"url"`
}

type RtorrentConfig struct {
	URL       string
	Transport http.RoundTripper
}

type Rtorrent struct {
	client *xmlrpc.Client
}

// Creates a new instance of Rtorrent client
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

// Lists available XMLRPC methods
func (rt *Rtorrent) ListMethods() ([]string, error) {
	var result []string
	err := rt.client.Call("system.listMethods", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Load and start a torrent
func (rt *Rtorrent) LoadRawStart(file []byte) error {
	base64 := base64.StdEncoding.EncodeToString(file)

	err := rt.client.Call("load.raw_start_verbose", []interface{}{"", xmlrpc.Base64(base64)}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Stop torrent with the specified hash
func (rt *Rtorrent) Stop(hash string) error {
	err := rt.client.Call("d.stop", hash, nil)
	if err != nil {
		return err
	}
	return nil
}

// Start torrent with the specified hash
func (rt *Rtorrent) Start(hash string) error {
	err := rt.client.Call("d.start", hash, nil)
	if err != nil {
		return err
	}
	return nil
}

func (rt *Rtorrent) DMulticall(view string, args interface{}) ([]Torrent, error) {
	var result interface{}
	err := rt.client.Call("d.multicall2", args, &result)
	if err != nil {
		return nil, err
	}

	torrents := multicallTags[Torrent](result, args)
	return torrents, nil
}

func (rt *Rtorrent) FMulticall(args interface{}) ([]File, error) {
	var result interface{}
	err := rt.client.Call("f.multicall", args, &result)
	if err != nil {
		return nil, err
	}

	files := multicallTags[File](result, args)
	return files, nil
}

func (rt *Rtorrent) PMulticall(args interface{}) ([]Peer, error) {
	var result interface{}
	err := rt.client.Call("p.multicall", args, &result)
	if err != nil {
		return nil, err
	}

	peers := multicallTags[Peer](result, args)
	return peers, nil
}

func (rt *Rtorrent) TMulticall(args interface{}) ([]Tracker, error) {
	var result interface{}
	err := rt.client.Call("t.multicall", args, &result)
	if err != nil {
		return nil, err
	}

	trackers := multicallTags[Tracker](result, args)
	return trackers, nil
}

// Maps XMLRPC result to a struct using fields from args with reflection
func multicallTags[T File | Torrent | Peer | Tracker](result interface{}, args interface{}) []T {
	items := make([]T, 0)
	for _, outer := range result.([]interface{}) {
		item := new(T)
		for idx := 2; idx < len(args.([]interface{})); idx++ {
			ref := outer.([]interface{})[idx-2]
			fname := args.([]interface{})[idx].(string)
			vo := reflect.ValueOf(item)
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
		items = append(items, *item)
	}
	return items
}
