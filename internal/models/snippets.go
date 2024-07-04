package models

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetForm struct {
	Title   string
	Content string
	Expires int
}

type SnippetModel struct {
	db *pgxpool.Pool
}

func NewSnippetModel(db *pgxpool.Pool) *SnippetModel {
	return &SnippetModel{db: db}
}

func (m *SnippetModel) NewSnippetForm(title string, content string, expires int) *SnippetForm {
	return &SnippetForm{
		Title:   title,
		Content: content,
		Expires: expires,
	}
}

func (m *SnippetModel) Validate(f *SnippetForm) (*map[string]string, bool) {
	snippetErrors := make(map[string]string)

	if strings.TrimSpace(f.Title) == "" {
		snippetErrors["titleError"] = "Title cannot be empty"
	} else if utf8.RuneCountInString(f.Title) > 100 {
		snippetErrors["titleError"] = "Title cannot be more than 100 characters"
	}

	if strings.TrimSpace(f.Content) == "" {
		snippetErrors["contentError"] = "Content cannot be empty"
	}

	if f.Expires < 0 {
		snippetErrors["expiresError"] = "Expiration cannot be negative"
	}

	if len(snippetErrors) > 0 {
		return &snippetErrors, false
	}

	return &snippetErrors, true
}

func (m *SnippetModel) Insert(f *SnippetForm) (int, error) {
	stmt := "INSERT INTO snippets(title, content, created, expires) VALUES($1, $2, now(), now() + MAKE_INTERVAL(DAYS => $3)) RETURNING id;"

	var id int
	err := m.db.QueryRow(context.Background(), stmt, f.Title, f.Content, f.Expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > now() AND id = $1;"
	var s = &Snippet{}
	err := m.db.QueryRow(context.Background(), stmt, id).Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > now() ORDER BY id DESC LIMIT 10;"

	rows, err := m.db.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
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
