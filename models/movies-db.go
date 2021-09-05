package models

import (
	"context"
	"database/sql"
	"time"
)

type DBmodel struct {
	DB *sql.DB
}

// Get returns one movie and an error, if any
func (m *DBmodel) Get(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at 
	from movies where id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)

	var movie Movie

	err := row.Scan(
		&movie.Id,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.Rating,
		&movie.MPAARating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	//genres if any
	query = `Select mg.id, mg.movie_id, mg.genre_id, g.genre_name
			from movies_genres mg
			left join genres g on g.id = mg.genre_id
			where mg.movie_id = $1`

	rows, _ := m.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	var genres []MovieGenre
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.Id,
			&mg.MovieId,
			&mg.GenreId,
			&mg.Genre,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, mg)
	}
	movie.MovieGenre = genres

	return &movie, nil
}

// All returns all movies and an error, if any
func (m *DBmodel) All() ([]*Movie, error) {

	return nil, nil
}
