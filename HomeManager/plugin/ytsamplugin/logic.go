package ytsamplugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/nielsvanm/homemanager/database"
	"github.com/nielsvanm/homemanager/tools/log"
)

/* TODO:
- Update create queries to break upon duplicates
- Add iterative loop to GetMovies to actually get all the movies
- Create Views
*/

// BaseURL is the api endpoint url
var BaseURL = "https://yts.am/api/v2/"

// RequestLimit is the max amount of movies we should request
var RequestLimit = 50

var totalMovieCount = 0
var currentPage = 1

type ResponseData struct {
	Status        string           `json:"status,omitempty"`
	StatusMessage string           `json:"status_message,omitempty"`
	Data          ResponseSettings `json:"data,omitempty"`
}

type ResponseSettings struct {
	MovieCount int     `json:"movie_count,omitempty"`
	PageNumber int     `json:"page_number,omitempty"`
	Movies     []Movie `json:"movies,omitempty"`
}

// GetMovies is the "main" function of the plugin
func GetMovies() []database.BatchQuery {
	// Create Queries
	genreBatch := database.BatchQuery{}
	genreBatch.Query = `
	INSERT INTO ytsamplugin_genre (name) 
	VALUES ($1) ON CONFLICT DO NOTHING;`

	movieBatch := database.BatchQuery{}
	movieBatch.Query = `
	INSERT INTO ytsamplugin_movie (title, year, rating, length, description, cover_image)
	VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING;`

	movieGenreBatch := database.BatchQuery{}
	movieGenreBatch.Query = `
	INSERT INTO ytsamplugin_movie_genre (movie_id, genre_id)
	VALUES (
		(SELECT id FROM ytsamplugin_movie
		WHERE title = $1),
		(SELECT id FROM ytsamplugin_genre
		WHERE name = $2)
		)
	ON CONFLICT DO NOTHING;`

	torrentBatch := database.BatchQuery{}
	torrentBatch.Query = `
	INSERT INTO ytsamplugin_torrent (quality, type, size, url, movie)
	VALUES ($1, $2, $3, $4, (
		SELECT id FROM ytsamplugin_movie
		WHERE title = $5
	))
	ON CONFLICT DO NOTHING;`

	// Repeat request until we run out of movies
	for {
		log.Info("YTSAMPlugin", "Requesting page "+strconv.Itoa(currentPage))

		resp := QueryYTS(currentPage)
		currentPage++

		if len(resp.Data.Movies) == 0 {
			break
		}

		// Parse data to queries
		for _, movie := range resp.Data.Movies {
			movieBatch.AddValues(
				movie.Title,
				movie.Year,
				movie.Rating,
				movie.Length,
				movie.Description,
				movie.CoverImage,
			)

			for _, genre := range movie.Genres {
				genreBatch.AddValues(
					genre,
				)

				movieGenreBatch.AddValues(
					movie.Title,
					genre,
				)
			}

			for _, torrent := range movie.Torrents {
				torrentBatch.AddValues(
					torrent.Quality,
					torrent.Type,
					torrent.Size,
					torrent.URL,
					movie.Title,
				)
			}
			log.Info("YTSAMPlugin", "Succesfully created queries for", movie.Title)
		}
	}

	// Create pluginresult
	return []database.BatchQuery{
		genreBatch,
		movieBatch,
		movieGenreBatch,
		torrentBatch,
	}
}

// QueryYTS queries the YTS.AM API at the provided page
func QueryYTS(page int) *ResponseData {
	templateURL := "%slist_movies.json?limit=%d&page=%d"
	URL := fmt.Sprintf(templateURL, BaseURL, RequestLimit, page)

	// Make the request and read the body
	response, err := http.Get(URL)
	if err != nil {
		log.Warn("YTSAMPlugin", err.Error())
	}

	defer response.Body.Close()

	log.Info("YTSAMPlugin", "Retrieved info from YTS.AM")

	blob, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Warn("YTSAMPlugin", err.Error())
	}

	// Parse the response json
	resp := ResponseData{}
	err = json.Unmarshal(blob, &resp)
	if err != nil {
		log.Warn("YTSAMPlugin", err.Error())
		return nil
	}

	if resp.Status != "ok" {
		log.Err("YTSAMPlugin", "Query failed: "+resp.StatusMessage)
		return nil
	}

	return &resp
}
