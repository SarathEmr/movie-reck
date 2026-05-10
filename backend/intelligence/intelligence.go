package intelligence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"movie-reck/model"
	"net/http"

	"encoding/csv"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func RecommendMoviesUsingAI(req model.RecommendationRequest) ([]model.Movie, error) {

	log.Print("--- in RecommendMoviesUsingAI\n\n") //

	apiKey := viper.GetString("GEMINI_API_KEY")
	if apiKey == "" {
		return []model.Movie{}, errors.New("GEMINI_API_KEY not set")
	}
	modelName := viper.GetString("GEMINI_MODEL_NAME")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", modelName, apiKey)

	genresDefault := []string{string(model.GenreSciFi), string(model.GenreDrama)}
	var genres []string
	if len(req.Genres) > 0 {
		for _, g := range req.Genres {
			genres = append(genres, string(g))
		}
	} else {
		genres = genresDefault
	}

	prompt := "Give me a list of some top-rated movies as per the below conditions:\n" +
		"Consider all (or at least one) from the following genres when suggesting the movies.\n" +
		"Language(s): English\n" +
		"Year range: 2010 to 2026\n" +
		"Movie count: 3\n" +
		"Genres: %s\n" +
		"Response should be in CSV format; with rows as follows: name, year, genre(s), imdb rating, why it fits.\n" +
		"The response string should not contain unnecessary format specifiers.\n" +
		"In each record other than the header, only the name, genre(s) and `why it fits` section should be wrapped in quotes."

	prompt = fmt.Sprintf(prompt, strings.Join(genres, ", ")) //

	reqBody := GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
	}

	b, _ := json.Marshal(reqBody)
	log.Printf("--- llmReq body = [%s], url = [%s]\n\n", b, url) //

	llmReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	llmReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(llmReq)
	if err != nil {
		return []model.Movie{}, fmt.Errorf("LLM request failed: %s", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("----- raw resp body is %s", string(body)) //

	var gr GeminiResponse
	if err := json.Unmarshal(body, &gr); err != nil {
		return []model.Movie{}, fmt.Errorf("failed to unmarshal LLM response: %s", err)
	}
	log.Printf("--- LLM response: %s\n\n", gr)

	if len(gr.Candidates) == 0 {
		return []model.Movie{}, errors.New("no response from LLM")
	}

	csvResponse := gr.Candidates[0].Content.Parts[0].Text
	fmt.Println("-- csvInput is ", csvResponse) //

	movies, err := ParseMoviesCSV(csvResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response csv: %s", err)
	}
	log.Print("--- movies parsed from LLM response: ", movies) //

	return movies, nil
}

// export GEMINI_API_KEY="your-api-key"

// ParseMoviesCSV parses the CSV data to Movie type in Go.
func ParseMoviesCSV(csvData string) ([]model.Movie, error) {
	reader := csv.NewReader(strings.NewReader(csvData))
	reader.TrimLeadingSpace = true

	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var movies []model.Movie

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// CSV columns:
		// name, year, genre(s), imdb rating, why it fits
		name := record[0]

		yearInt, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}

		imdb, err := strconv.ParseFloat(record[3], 32)
		if err != nil {
			return nil, err
		}

		movie := model.Movie{
			ID:         generateID(name, yearInt),
			Name:       name,
			Year:       int16(yearInt),
			Language:   "English",
			ImdbRating: float32(imdb),
			Genres:     parseGenres(record[2]),
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

// ----- Helpers -----

func parseGenres(raw string) []model.Genre {
	parts := strings.Split(raw, ",")
	var genres []model.Genre

	for _, p := range parts {
		switch strings.ToLower(strings.TrimSpace(p)) {
		case "sci-fi", "science fiction":
			genres = append(genres, model.GenreSciFi)
		case "drama":
			genres = append(genres, model.GenreDrama)
		case "romcom":
			genres = append(genres, model.GenreRomcom)
		case "horror":
			genres = append(genres, model.GenreHorror)
		case "horror-comedy":
			genres = append(genres, model.GenreHorrorComedy)
		case "thriller":
			genres = append(genres, model.GenreThriller)

		case "romance":
			genres = append(genres, model.GenreRomance)
		case "comedy":
			genres = append(genres, model.GenreComedy)
		case "music":
			genres = append(genres, model.GenreMusic)
		}
	}

	return genres
}

func generateID(name string, year int) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	return fmt.Sprintf("%s-%d", slug, year)
}
