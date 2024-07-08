package models

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserDao struct {
	Id       int
	Name     string
	Email    string
	Password string
	Created  time.Time
}

type User struct {
	Id      int
	Name    string
	Email   string
	Created time.Time
}

type UserSignupForm struct {
	Name     string
	Email    string
	Password string
}

type UserLoginForm struct {
	Email    string
	Password string
}

type UserModel struct {
	db *pgxpool.Pool
}

func NewUserModel(db *pgxpool.Pool) *UserModel {
	return &UserModel{db: db}
}

func (m *UserModel) NewUserSignupForm(name, email, password string) *UserSignupForm {
	return &UserSignupForm{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func (m *UserModel) NewUserLoginForm(email, password string) *UserLoginForm {
	return &UserLoginForm{
		Email:    email,
		Password: password,
	}
}

func (m *UserModel) Validate(f *UserSignupForm) (*map[string]string, bool) {
	userErrors := make(map[string]string)

	if strings.TrimSpace(f.Name) == "" {
		userErrors["nameError"] = "Name cannot be empty"
	} else if utf8.RuneCountInString(f.Name) > 100 {
		userErrors["nameError"] = "Name cannot be more than 100 characters"
	}

	if strings.TrimSpace(f.Password) == "" {
		userErrors["passwordError"] = "Password cannot be empty"
	} else if utf8.RuneCountInString(f.Password) < 8 {
		userErrors["passwordError"] = "Password must be at least 8 characters"
	}

	if len(userErrors) > 0 {
		return &userErrors, false
	}

	return &userErrors, true
}

func (m *UserModel) Signup(f *UserSignupForm) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(f.Password), 12)
	if err != nil {
		return 0, err
	}
	stmt := "INSERT INTO users(name, email, password, created) VALUES($1, $2, $3, now()) RETURNING id;"

	var id int
	err = m.db.QueryRow(context.Background(), stmt, f.Name, f.Email, hash).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *UserModel) Login(f *UserLoginForm) (*User, error) {
	stmt := "SELECT id, name, email, created, password FROM users WHERE email = $1;"

	var u = &User{}
	var hashedPassword string
	err := m.db.QueryRow(context.Background(), stmt, f.Email).Scan(&u.Id, &u.Name, &u.Email, &u.Created, &hashedPassword)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(f.Password)); err != nil {

		return nil, errors.New("incorrect password")
	}

	return u, nil
}

func (m *UserModel) GetById(id int) (*User, error) {
	stmt := "SELECT id, name, email, created FROM users WHERE id = $1;"

	var u = &User{}
	if err := m.db.QueryRow(context.Background(), stmt, id).Scan(&u.Id, &u.Name, &u.Email, &u.Created); err != nil {
		return nil, err
	}

	return u, nil
}
