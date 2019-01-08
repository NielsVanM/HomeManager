package torrentplugin

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"nvmtech.nl/homemanager/tools"

	"github.com/lnguyen/go-transmission/transmission"
	"nvmtech.nl/homemanager/frame"
	"nvmtech.nl/homemanager/tools/log"
)

var tmClient = transmission.New("http://localhost:9091", "", "")
var TorrentFolder = "./__data/torrentplugin/torrents/"

var APIEndpoints = []*frame.Endpoint{
	frame.NewEndpoint("/add/", APIAddTorrentView),
}
var ViewEndpoints = []*frame.Endpoint{
	frame.NewEndpoint("/", DashboardView),
}

func DashboardView(w http.ResponseWriter, r *http.Request) {
	page := frame.NewPage([]string{"base.html", "plugins/torrentplugin/dashboard.html"})

	torrents, err := tmClient.GetTorrents()
	if err != nil {
		log.Warn("TorrentPlugin", "Failed to get torrents from transmission", err.Error())
		return
	}

	for i := range torrents {
		// Fix percent done
		torrents[i].PercentDone = torrents[i].PercentDone * 100

		// Fix Download speed
		torrents[i].DownloadDir = tools.ByteCountDecimal(int64(torrents[i].RateDownload))
	}

	page.AddContext("torrents", torrents)

	page.Render(w)
}

// APIAddTorrentView allows an external program or plugin to send a torrent to
// download
func APIAddTorrentView(w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("torrentfile")
	if err != nil {
		log.Warn("TorrentPlugin", "Failed to retrieve file", err.Error())
		return
	}
	defer file.Close()

	// Add .torrent extension
	if !strings.HasSuffix(handle.Filename, ".torrent") {
		handle.Filename += ".torrent"
	}

	// Create/open file and copy contents from received file into our file
	f, err := os.OpenFile(TorrentFolder+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Warn("TorrentPlugin", "Failed to open or create file at", TorrentFolder+handle.Filename, err.Error())
		return
	}

	_, err = io.Copy(f, file)
	if err != nil {
		log.Warn("TorrentPlugin", "Failed to copy data", err.Error())
	}
	f.Close()

	// Add torrent transmission
	ta, err := tmClient.AddTorrentByFilename(TorrentFolder+handle.Filename, "./__data/torrentplugin/downloads/")
	if err != nil {
		log.Warn("TorrentPlugin", "Failed to add torrent", err.Error())
		return
	}

	log.Info("TorrentPlugin", ta.Name, string(ta.ID), ta.HashString)

	// Respond
	resp, err := json.Marshal(`{"status": 200, "status_text": "Torrent Succesfully added"}`)
	if err != nil {
		log.Warn("TorrentPlugin", "Failed to Marshal json response")
		return
	}

	w.Write(resp)
}
