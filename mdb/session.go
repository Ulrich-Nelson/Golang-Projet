package mdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const LoginFilename = "login.json"

type Session struct {
	db *sql.DB
}

func NewSession() *Session {
	login := readLogin()
	db := connect(login)
	return &Session{db: db}
}

func (s *Session) Close() {
	defer s.db.Close()
}

type Login struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func readLogin() Login {
	file, err := os.Open(LoginFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var login Login
	err = json.Unmarshal(data, &login)
	if err != nil {
		panic(err)
	}
	return login
}

func connect(login Login) *sql.DB {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		login.Host, login.Port, login.User, login.Password, login.DBName,
	)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}
	return db
}
