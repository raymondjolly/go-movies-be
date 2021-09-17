package models

import (
	"context"
	"database/sql"
	"fmt"
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

	genres := make(map[int]string)
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
		genres[mg.Id] = mg.Genre
	}
	movie.MovieGenre = genres

	return &movie, nil
}

// All returns all movies and an error, if any. The function can (but doesn't have to) take a parameter of a
// genre index
func (m *DBmodel) All(genre ...int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	query := fmt.Sprintf(`select id, title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at 
	from movies %s order by title`, where)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie

	for rows.Next() {
		var movie Movie
		err := rows.Scan(
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
		genreQuery := `Select mg.id, mg.movie_id, mg.genre_id, g.genre_name
			from movies_genres mg
			left join genres g on g.id = mg.genre_id
			where mg.movie_id = $1`

		genreRows, _ := m.DB.QueryContext(ctx, genreQuery, movie.Id)

		genres := make(map[int]string)
		for genreRows.Next() {
			var mg MovieGenre
			err := genreRows.Scan(
				&mg.Id,
				&mg.MovieId,
				&mg.GenreId,
				&mg.Genre,
			)
			if err != nil {
				return nil, err
			}
			genres[mg.Id] = mg.Genre
		}
		genreRows.Close()
		movie.MovieGenre = genres
		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *DBmodel) GenresAll() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `select id, genre_name, created_at, updated_at from genres order by genre_name`

	rows, err := m.DB.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var genres []*Genre
	for rows.Next() {
		var g Genre
		err = rows.Scan(
			&g.Id,
			&g.GenreName,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}
	return genres, nil
}

func (m *DBmodel) InsertMovie(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into movies(title, description, year, release_date, runtime, rating, mpaa_rating, created_at,
		updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`

	_, err := m.DB.ExecContext(ctx, stmt, movie.Title, movie.Description, movie.Year, movie.ReleaseDate, movie.RunTime,
		movie.Rating, movie.MPAARating, movie.CreatedAt, movie.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}
