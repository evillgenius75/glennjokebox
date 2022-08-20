package models

import (
	"database/sql"
	"errors"
	"time"
)

type Joke struct {
	ID        int
	UUID      string
	Joke      string
	Explicit  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JokeModel struct {
	DB *sql.DB
}

func (m *JokeModel) Insert(uuid string, joke string, explicit int) (int, error) {
	stmt := `INSERT INTO joke_results (uuid, joke, explicit, created_at, updated_at)
VALUES(?,?,?,UTC_TIMESTAMP(),UTC_TIMESTAMP())`
	result, err := m.DB.Exec(stmt, uuid, joke, explicit)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *JokeModel) Get(id int) (*Joke, error) {
	stmt := `SELECT id, uuid, joke, explicit, created_at, updated_at FROM joke_results
WHERE id=?`
	row := m.DB.QueryRow(stmt, id)
	j := &Joke{}

	err := row.Scan(&j.ID, &j.UUID, &j.Joke, &j.Explicit, &j.CreatedAt, &j.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return j, nil
}

func (m *JokeModel) Latest() ([]*Joke, error) {
	stmt := `SELECT id, uuid, joke, explicit, created_at, updated_at FROM joke_results
ORDER BY RAND() LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jokes := []*Joke{}
	for rows.Next() {
		j := &Joke{}
		err := rows.Scan(&j.ID, &j.UUID, &j.Joke, &j.Explicit, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jokes = append(jokes, j)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jokes, nil
}
