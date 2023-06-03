package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"reflect"

	"github.com/kolo/xmlrpc"
)

type Torrent struct {
	Hash           string `rtw:"d.hash=" json:"hash"`
	Name           string `rtw:"d.name=" json:"name"`
	SizeBytes      int64  `rtw:"d.size_bytes=" json:"size_bytes"`
	CompletedBytes int64  `rtw:"d.completed_bytes=" json:"completed_bytes"`
	UploadRate     int64  `rtw:"d.up.rate=" json:"upload_rate"`
	UploadTotal    int64  `rtw:"d.up.total=" json:"upload_total"`
	DownloadRate   int64  `rtw:"d.down.rate=" json:"download_rate"`
	DownloadTotal  int64  `rtw:"d.down.total=" json:"download_total"`
	Message        string `rtw:"d.message=" json:"message"`
	IsActive       int64  `rtw:"d.is_active=" json:"is_active"`
	IsOpen         int64  `rtw:"d.is_open=" json:"is_open"`
	IsHashing      int64  `rtw:"d.is_hash_checking=" json:"is_hashing"`
	Leechers       int64  `rtw:"d.peers_accounted=" json:"leechers"`
	Seeders        int64  `rtw:"d.peers_complete=" json:"seeders"`
	State          int64  `rtw:"d.state=" json:"state"`
	StateChanged   int64  `rtw:"d.state_changed=" json:"state_changed"`
	StateCounter   int64  `rtw:"d.state_counter=" json:"state_counter"`
	Priority       int64  `rtw:"d.priority=" json:"priority"`
	Custom1        string `rtw:"d.custom1=" json:"custom1"`
	Custom2        string `rtw:"d.custom2=" json:"custom2"`
	Custom3        string `rtw:"d.custom3=" json:"custom3"`
	Custom4        string `rtw:"d.custom4=" json:"custom4"`
	Custom5        string `rtw:"d.custom5=" json:"custom5"`
}

type File struct {
	Path            string `rtw:"f.path=" json:"path"`
	Size            int64  `rtw:"f.size_bytes=" json:"size"`
	SizeChunks      int64  `rtw:"f.size_chunks=" json:"size_chunks"`
	CompletedChunks int64  `rtw:"f.completed_chunks=" json:"completed_chunks"`
	FrozenPath      string `rtw:"f.frozen_path=" json:"frozen_path"`
	Priority        int64  `rtw:"f.priority=" json:"priority"`
	IsCreated       int64  `rtw:"f.is_created=" json:"is_created"`
	IsOpen          int64  `rtw:"f.is_open=" json:"is_open"`
}

type Peer struct {
	PeerID           string `rtw:"p.id=" json:"peer_id"`
	Address          string `rtw:"p.address=" json:"address"`
	Port             int64  `rtw:"p.port=" json:"port"`
	Banned           int64  `rtw:"p.banned=" json:"banned"`
	ClientVersion    string `rtw:"p.client_version=" json:"client_version"`
	CompletedPercent int64  `rtw:"p.completed_percent=" json:"completed_percent"`
	IsEncrypted      int64  `rtw:"p.is_encrypted=" json:"is_encrypted"`
	IsIncoming       int64  `rtw:"p.is_incoming=" json:"is_incoming"`
	IsObfuscated     int64  `rtw:"p.is_obfuscated=" json:"is_obfuscated"`
	PeerRate         int64  `rtw:"p.peer_rate=" json:"peer_rate"`
	PeerTotal        int64  `rtw:"p.peer_total=" json:"peer_total"`
	UploadRate       int64  `rtw:"p.up_rate=" json:"up_rate"`
	UploadTotal      int64  `rtw:"p.up_total=" json:"up_total"`
}

type Tracker struct {
	TrackerID        string `rtw:"t.id=" json:"tracker_id"`
	ActivityTimeLast int64  `rtw:"t.activity_time_last=" json:"activity_time_last"`
	ActivityTimeNext int64  `rtw:"t.activity_time_next=" json:"activity_time_next"`
	CanScrape        int64  `rtw:"t.can_scrape=" json:"can_scrape"`
	IsUsable         int64  `rtw:"t.is_usable=" json:"t.is_usable"`
	IsEnabled        int64  `rtw:"t.is_enabled=" json:"is_enabled"`
	FailedCounter    int64  `rtw:"t.failed_counter=" json:"failed_counter"`
	FailedTimeLast   int64  `rtw:"t.failed_time_last=" json:"failed_time_last"`
	FailedTimeNext   int64  `rtw:"t.failed_time_next=" json:"failed_time_next"`
	IsBusy           int64  `rtw:"t.is_busy=" json:"is_busy"`
	IsOpen           int64  `rtw:"t.is_open=" json:"is_open"`
	Type             int64  `rtw:"t.type=" json:"type"`
	URL              string `rtw:"t.url=" json:"url"`
}

type System struct {
	APIVersion     string `rtw:"system.api_version" json:"api_version"`
	ClientVersion  string `rtw:"system.client_version" json:"client_version"`
	LibraryVersion string `rtw:"system.library_version" json:"library_version"`
}

type SystemCall struct {
	MethodName string      `xmlrpc:"methodName" json:"method_name"`
	Params     interface{} `xmlrpc:"params" json:"params"`
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

func (rt *Rtorrent) SystemMulticall(args interface{}) (System, error) {
	var result interface{}
	err := rt.client.Call("system.multicall", args, &result)
	if err != nil {
		return System{}, err
	}

	system := systemTags(result, args)
	return system, nil
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
				if fname == field.Tag.Get("rtw") {
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

func systemTags(result interface{}, args interface{}) System {
	system := &System{}

	// todo: do work

	fmt.Printf("%v %v", result, args)

	return *system
}
