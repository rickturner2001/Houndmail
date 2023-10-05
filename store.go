package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type MySqlStore struct {
	Db *sql.DB
}

func NewMySqlStore() *MySqlStore {
	db, err := sql.Open("mysql", "admin:password@/development")
	if err != nil {
		log.Fatalf("Could not connect to db: %+v", err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatalf("Could not ping to db: %+v", err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	return &MySqlStore{
		Db: db,
	}
}

func (s MySqlStore) Init() {
	createTableStmt := `
	CREATE TABLE IF NOT EXISTS users(
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL UNIQUE
	)
	`
	_, err := s.Db.Exec(createTableStmt)
	if err != nil {
		log.Fatalf("Could not execute SQL statement:\n%s\n[ERROR]: %+v", createTableStmt, err)
	}

	log.Println("Successfully created users table")
}

func (s MySqlStore) RegisterUser(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Could not create user %s: %s", user.Username, err)
	}

	_, err = s.Db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, hashedPassword)
	return err
}

func (s MySqlStore) AuthenticateUser(user User) bool {
	var hashedPassword string
	err := s.Db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&hashedPassword)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	return err == nil
}
