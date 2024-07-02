package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	db *pgxpool.Pool
}

func NewSnippetModel(db *pgxpool.Pool) *SnippetModel {
	return &SnippetModel{db: db}
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := "INSERT INTO snippets(title, content, created, expires) VALUES($1, $2, now(), now() + MAKE_INTERVAL(DAYS => $3)) RETURNING id;"

	var id int
	err := m.db.QueryRow(context.Background(), stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := "SELECT title, content, created, expires FROM snippets WHERE expires > now() AND id = $1;"
	var s = &Snippet{}
	err := m.db.QueryRow(context.Background(), stmt, id).Scan(&s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := "SELECT title, content, created, expires FROM snippets WHERE expires > now() ORDER BY id DESC LIMIT 10;"

	rows, err := m.db.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
