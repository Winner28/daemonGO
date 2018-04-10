package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	// psql
	_ "github.com/lib/pq"
)

// Handler struct represents connection to db
type Handler struct {
	Db     *sql.DB
	config Config
}

// Config represents config of db
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// GetConfig returns empty config
func GetConfig() Config {
	return Config{}
}

// GetConnection returns a connection
func GetConnection(config Config) (handler *Handler, err error) {
	if !validateConfig(config) {
		return nil, errors.New("Can`t open connection with empty fields")
	}

	handler = new(Handler)

	handler.config = config

	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.User, config.Password, config.Database, config.Host, config.Port)

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		log.Println("Error when trying to connect db")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Couldn`t establish a connection with the database")
		return nil, errors.New("Couldn`t establish a connection with the database")
	}
	handler.Db = db
	log.Println("Connection created")
	return handler, nil
}

// CloseConnection closes connection to DB
func (handler *Handler) CloseConnection() error {
	if handler == nil || handler.Db == nil {
		return nil
	}

	if err := handler.Db.Close(); err != nil {
		log.Println("Connection not closed")
		return errors.New("Error when trying to close connection")
	}

	log.Println("Connection closed")
	return nil
}

func validateConfig(config Config) bool {
	return !(config.Database == "" || config.Host == "" ||
		config.Password == "" || config.Port == "" ||
		config.User == "")
}
