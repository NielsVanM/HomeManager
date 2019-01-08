package ytsamplugin

import (
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nielsvanm/homemanager/frame"
)

// APIEndpoints List of endpoints for the API
var APIEndpoints = []*frame.Endpoint{}

// ViewEndpoints List of endpoints for the webapp
var ViewEndpoints = []*frame.Endpoint{
	frame.NewEndpoint("/", DashboardView),
	frame.NewEndpoint("/movie/", MovieOverviewView),
	frame.NewEndpoint("/search/title/", TitleSearch),
	frame.NewEndpoint("/view/{movieid}/", MovieView),
	frame.NewEndpoint("/torrent/{torrentid}/", DownloadTorrentView),
}

// DashboardView renders the dashboard template
func DashboardView(w http.ResponseWriter, r *http.Request) {
	dashTemplate := frame.NewPage([]string{"base.html", "plugins/ytsamplugin/dashboard.html"})

	dashTemplate.Render(w)
}

// MovieOverviewView renders the movies for the dashboard
func MovieOverviewView(w http.ResponseWriter, r *http.Request) {
	movieTemplate := frame.NewPage([]string{"plugins/ytsamplugin/movieoverview.html"})

	getDownloaded := false
	getDownloaded, err := strconv.ParseBool(r.URL.Query().Get("download"))
	if err != nil {
		getDownloaded = false
	}

	var page int64
	page, err = strconv.ParseInt(
		r.URL.Query().Get("page"),
		10,
		64,
	)

	if err != nil {
		page = 0
	}

	pageSize := 48

	// Get all movies from db
	movieList := GetAllMovies(pageSize, int(page)*pageSize, getDownloaded)
	movieTemplate.AddContext("movies", movieList)

	// Get all years
	// years := GetUniqueYears(getDownloaded)
	// movieTemplate.AddContext("years", years)

	// Get page count
	// pageCount := GetPageCount(int(page), pageSize, getDownloaded)
	// movieTemplate.AddContext("pages", pageCount)

	movieTemplate.Render(w)
}

func TitleSearch(w http.ResponseWriter, r *http.Request) {
	movieTemplate := frame.NewPage([]string{"plugins/ytsamplugin/movieoverview.html"})

	title := r.URL.Query().Get("title")
	if title == "" {
		MovieOverviewView(w, r)
		return
	}

	movies := GetMovieByTitle(title, 0, 50, false)

	movieTemplate.AddContext(
		"movies", movies,
	)

	movieTemplate.Render(w)
}

func MovieView(w http.ResponseWriter, r *http.Request) {
	page := frame.NewPage([]string{"base.html", "plugins/ytsamplugin/movieview.html"})

	movieID, _ := strconv.Atoi(mux.Vars(r)["movieid"])

	movie := GetSingleMovie(movieID)
	page.AddContext("movie", movie)

	genres := GetGenreByMovie(movie.ID)
	genreString := strings.Join(genres, "/")
	page.AddContext("genres", genreString)

	torrents := GetTorrentsByMovie(movie.ID)
	page.AddContext("torrents", torrents)

	page.Render(w)
}

// DownloadTorrentView downloads a torrent in the background
func DownloadTorrentView(w http.ResponseWriter, r *http.Request) {
	torrentID, _ := strconv.Atoi(mux.Vars(r)["torrentid"])

	torrent := GetTorrentByID(torrentID)

	torrentDir := "./__data/ytsamplugin/torrents/"

	tmp := exec.Command("wget", "--directory-prefix="+torrentDir, torrent.URL)
	tmp.Start()
}
