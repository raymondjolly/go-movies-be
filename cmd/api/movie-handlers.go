package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go-movies-be/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.Models.DB.Get(id)

	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.Models.DB.All()
	if err != nil {
		app.errorJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.Models.DB.GenresAll()
	if err != nil {
		app.errorJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, genres, "genres")
}

func (app *application) getMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	genreId, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		app.errorJSON(w, err)
	}
	movies, err := app.Models.DB.All(genreId)
	if err != nil {
		app.errorJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
}

type MoviePayload struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	RunTime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

func (app *application) editMovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	var movie models.Movie

	movie.Id, _ = strconv.Atoi(payload.Id)
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.Rating, _ = strconv.Atoi(payload.Rating)
	movie.ReleaseDate, _ = time.Parse("2006-01-02", payload.ReleaseDate)
	movie.RunTime, _ = strconv.Atoi(payload.RunTime)
	movie.Year = movie.ReleaseDate.Year()
	movie.MPAARating = payload.MPAARating
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	if movie.Id == 0 {
		err = app.Models.DB.InsertMovie(movie)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	ok := jsonResp{
		OK: true,
	}

	err = app.writeJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJSON(w, err)
	}
}

//TODO write a handler for a deleteMovie function

//TODO write a handler for a searchMovie function
