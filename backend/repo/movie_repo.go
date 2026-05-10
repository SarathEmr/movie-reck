package repo

import (
	"database/sql"
	"fmt"
	"movie-reck/model"

	"github.com/lib/pq"
)

type MovieRepo struct {
	db *sql.DB
}

func NewMovieRepo(db *sql.DB) *MovieRepo {
	return &MovieRepo{db: db}
}

func (r *MovieRepo) GetAllMovies() ([]model.Movie, error) {
	rows, err := r.db.Query("SELECT id, name, year, language, imdb_rating, genre FROM movie")
	if err != nil {
		return nil, fmt.Errorf("db query failed: %s", err)
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		var genres []string
		err := rows.Scan(&movie.ID, &movie.Name, &movie.Year, &movie.Language, &movie.ImdbRating, pq.Array(&genres))
		if err != nil {
			return nil, err
		}
		mappedGenres, err := mapGenresToModel(genres)
		if err != nil {
			return nil, fmt.Errorf("mapGenresToModel failed: %s", err)
		}
		movie.Genres = mappedGenres
		movies = append(movies, movie)
	}
	return movies, nil
}

func mapGenresToModel(genres []string) ([]model.Genre, error) {
	genresMapped := make([]model.Genre, len(genres))
	for _, g := range genres {
		switch g {
		case string(model.GenreSciFi), string(model.GenreDrama), string(model.GenreRomcom), string(model.GenreHorror),
			string(model.GenreHorrorComedy), string(model.GenreThriller):
			genresMapped = append(genresMapped, model.Genre(g))
		default:
			return nil, fmt.Errorf("invalid genre: %s", g)
		}
	}
	return genresMapped, nil
}
