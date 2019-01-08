package ytsamplugin

import (
	"github.com/nielsvanm/homemanager/database"
	"github.com/nielsvanm/homemanager/tools/log"
)

// Movie is a representation of the movie data
type Movie struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Year        int       `json:"year,omitempty"`
	Rating      float32   `json:"rating,omitempty"`
	Length      int       `json:"runtime,omitempty"`
	Description string    `json:"description_full,omitempty"`
	CoverImage  string    `json:"large_cover_image,omitempty"`
	Genres      []string  `json:"genres,omitempty"`
	Torrents    []Torrent `json:"torrents,omitempty"`
}

// Torrent is the representation of torrent data
type Torrent struct {
	ID      int    `json:"id,omitempty"`
	Quality string `json:"quality,omitempty"`
	Type    string `json:"type,omitempty"`
	Size    string `json:"size,omitempty"`
	URL     string `json:"url,omitempty"`
}

// GetAllMovies returns a list of movies based on the limit, offset and
// downloaded filters
func GetAllMovies(limit, offset int, downloaded bool) []Movie {
	movieRows := database.Database.Query(`
	SELECT id, title, cover_image, year, rating, length
	FROM ytsamplugin_movie
	WHERE downloaded = $1
	ORDER BY random()
	LIMIT $2
	OFFSET $3;`, downloaded, limit, offset)

	defer movieRows.Close()

	movieList := []Movie{}

	for movieRows.Next() {
		var newMovie = Movie{}
		err := movieRows.Scan(
			&newMovie.ID,
			&newMovie.Title,
			&newMovie.CoverImage,
			&newMovie.Year,
			&newMovie.Rating,
			&newMovie.Length,
		)

		if err != nil {
			log.Warn("YTSAMPlugin", "Failed to scan row, closing stream")
			movieRows.Close()
		}

		movieList = append(movieList, newMovie)
	}

	return movieList
}

func GetSingleMovie(id int) Movie {
	rows := database.Database.Query(`
	SELECT id, title, year, rating, length, description, cover_image
	FROM ytsamplugin_movie
	WHERE id = $1`, id)

	for rows.Next() {
		movie := Movie{}

		rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Year,
			&movie.Rating,
			&movie.Length,
			&movie.Description,
			&movie.CoverImage,
		)
		return movie
	}
	return Movie{}
}

// GetUniqueYears returns a list of years that we have movies for
func GetUniqueYears(downloaded bool) []int {
	// Get all unique years
	yearRows := database.Database.Query(`
	SELECT DISTINCT year FROM ytsamplugin_movie
	WHERE downloaded = $1
	ORDER BY year DESC;
	`, downloaded)
	years := []int{}
	for yearRows.Next() {
		var year int
		err := yearRows.Scan(
			&year,
		)

		if err != nil {
			log.Warn("YTSAMPlugin", "Failed to scan row, closing stream")
			yearRows.Close()
		}

		years = append(years, year)
	}

	return years
}

func GetPageCount(current_page, page_size int, downloaded bool) []int {
	// Get a count of the pages
	pageRows := database.Database.Query(`
	SELECT COUNT(id) FROM ytsamplugin_movie
	WHERE downloaded = $1;`, downloaded)

	defer pageRows.Close()

	totalCount := 0
	for pageRows.Next() {
		pageRows.Scan(
			&totalCount,
		)
	}

	// Calculate pages
	pageCount := totalCount / page_size

	pages := []int{}
	for i := current_page - 10; i < current_page+10; i++ {
		if i < 0 || i > pageCount {
			continue
		}
		pages = append(pages, i+1)
	}

	return pages
}

// GetMovieByTitle returns the search results based on the title field
func GetMovieByTitle(title string, page, limit int, downloaded bool) []Movie {

	movieRows := database.Database.Query(`
	SELECT id, title, cover_image, year, rating, length
	FROM ytsamplugin_movie
	WHERE downloaded = $1 AND
	to_tsvector('english', title) @@ to_tsquery('english', $2);`,
		downloaded, string(title))

	defer movieRows.Close()

	movieList := []Movie{}

	for movieRows.Next() {
		var newMovie = Movie{}
		err := movieRows.Scan(
			&newMovie.ID,
			&newMovie.Title,
			&newMovie.CoverImage,
			&newMovie.Year,
			&newMovie.Rating,
			&newMovie.Length,
		)

		if err != nil {
			log.Warn("YTSAMPlugin", "Failed to scan row, closing stream")
			movieRows.Close()
		}

		movieList = append(movieList, newMovie)
	}

	return movieList
}

// GetGenreByMovie returns a list of genres for the provided movie
func GetGenreByMovie(movieID int) []string {
	genreRows := database.Database.Query(`
	SELECT name FROM ytsamplugin_genre
	WHERE id IN (
		SELECT genre_id FROM ytsamplugin_movie_genre
		WHERE movie_id = $1
	);`, movieID)

	genres := []string{}
	for genreRows.Next() {
		var genre string

		err := genreRows.Scan(
			&genre,
		)

		if err != nil {
			log.Warn("YTSAMPlugin", "Failed to parse genre rows")
		}

		genres = append(genres, genre)
	}

	return genres
}

// GetTorrentsByMovie returns a list of torrents associated with the provided
// movie
func GetTorrentsByMovie(movieID int) []Torrent {
	torrentRows := database.Database.Query(`
	SELECT id, quality, type, size, url FROM ytsamplugin_torrent
	WHERE movie = $1`, movieID)

	torrents := []Torrent{}
	for torrentRows.Next() {
		torrent := Torrent{}

		err := torrentRows.Scan(
			&torrent.ID,
			&torrent.Quality,
			&torrent.Type,
			&torrent.Size,
			&torrent.URL,
		)

		if err != nil {
			log.Warn("YTSAMPlugin", "Failed to parse torrent row")
		}

		torrents = append(torrents, torrent)
	}

	return torrents
}

// GetTorrentByID retrieves a torrent from the database specified by the id
func GetTorrentByID(torrentID int) *Torrent {
	torrentRows := database.Database.Query(`
	SELECT id, quality, type, size, url FROM ytsamplugin_torrent
	WHERE id = $1`, torrentID)

	for torrentRows.Next() {
		torrent := Torrent{}

		err := torrentRows.Scan(
			&torrent.ID,
			&torrent.Quality,
			&torrent.Type,
			&torrent.Size,
			&torrent.URL,
		)

		if err != nil {
			log.Warn("YTSAMPlugin", "Failed to parse torrent row")
		}

		return &torrent
	}

	return nil
}
