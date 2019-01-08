package ytsamplugin

var SetupDB = []string{
	`CREATE TABLE IF NOT EXISTS %s_genre (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);`,
	`CREATE TABLE IF NOT EXISTS %s_movie (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL UNIQUE,
		year INT,
		rating FLOAT,
		length INT,
		description TEXT,
		cover_image TEXT UNIQUE,
		downloaded BOOLEAN
	);`,
	`CREATE TABLE IF NOT EXISTS %s_torrent (
		id SERIAL PRIMARY KEY,
		quality TEXT,
		type TEXT,
		size TEXT,
		url text,
		movie INT REFERENCES %s_movie(id) ON UPDATE CASCADE,
		UNIQUE (movie, quality)
	);`,
	`CREATE TABLE IF NOT EXISTS %s_movie_genre (
		id SERIAL PRIMARY KEY,
		movie_id INT REFERENCES %s_movie(id) ON UPDATE CASCADE ON DELETE CASCADE,
		genre_id INT REFERENCES %s_genre(id) ON UPDATE CASCADE,
		UNIQUE (movie_id, genre_id)
	);`,
}

var Tables = []string{
	`%s_movie`,
	`%s_genre`,
	`%s_torrent`,
	`%s_movie_genre`,
}
