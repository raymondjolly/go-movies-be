package main

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"go-movies-be/models"
	"net/http"
	"strconv"
	"time"
)

func (app *application) getOneMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Println(errors.New("invalid id parameter"))
	}

	movie := models.Movie{
		Id:         id,
		Title:      "Dumb and Dumber",
		Year:       1994,
		MPAARating: "PG-13",
		RunTime:    153,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		app.logger.Println(err)
	}

}

func (app *application) getMovies(w http.ResponseWriter, r *http.Request) {

}
