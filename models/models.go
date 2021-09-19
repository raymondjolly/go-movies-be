package models

import (
	"database/sql"
	"time"
)

//Models is the wrapper for database
type Models struct {
	DB DBmodel
}

//NewModels returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBmodel{
			DB: db,
		},
	}
}

//Movie is the type for movies
type Movie struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Year        int            `json:"year"`
	ReleaseDate time.Time      `json:"release_date"`
	RunTime     int            `json:"runtime"`
	Rating      int            `json:"rating"`
	MPAARating  string         `json:"mpaa_rating"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	MovieGenre  map[int]string `json:"genres"`
}

//Genre is the type for genre
type Genre struct {
	Id        int       `json:"id"`
	GenreName string    `json:"genre_name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

//MovieGenre is the type for movie genre
type MovieGenre struct {
	Id        int       `json:"-"`
	MovieId   int       `json:"-"`
	GenreId   int       `json:"-"`
	Genre     string    `json:"genre"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

//User is the type for users
type User struct {
	Id       int
	Email    string
	Password string
}
