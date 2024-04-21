package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email    string
	Name     string
	Password []byte
}

type Service interface {
	AddUser(string, string, string) error
	GetUserByEmail(string) (User, error)
}

type service struct {
	db *sql.DB
}

var dburl = os.Getenv("DB_URL")

func New() Service {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		email TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	s := &service{db: db}
	return s
}

func (s *service) AddUser(email, name, password string) error {
	stmt, err := s.db.Prepare("INSERT INTO users (email, name, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(email, name, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (s *service) GetUserByEmail(email string) (User, error) {
	var user User

	row := s.db.QueryRow("SELECT * FROM users WHERE email = ?", email)

	err := row.Scan(&user.Email, &user.Name, &user.Password)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
