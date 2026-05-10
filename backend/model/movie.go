package model

type Genre string

const (
	GenreSciFi        Genre = "sci-fi"
	GenreDrama        Genre = "drama"
	GenreRomcom       Genre = "romcom"
	GenreHorror       Genre = "horror"
	GenreHorrorComedy Genre = "horror-comedy"
	GenreThriller     Genre = "thriller"

	GenreRomance Genre = "romance"
	GenreComedy  Genre = "comedy"
	GenreMusic   Genre = "music"
)

type Movie struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Year       int16   `json:"year"`
	Language   string  `json:"language"`
	ImdbRating float32 `json:"imdb_rating"`
	Genres     []Genre `json:"genre"`
}

type RecommendationRequest struct {
	Genres   []Genre `json:"genres" form:"genres"`
	Language string  `json:"language" form:"language"`
}
