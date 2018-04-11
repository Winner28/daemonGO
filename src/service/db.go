package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

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
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// GetEmptyConfig returns empty config
func GetEmptyConfig() Config {
	return Config{}
}

// GetDefaultConfig returns defaul config
func GetDefaultConfig() Config {
	config := Config{}
	config.User = "super_cool"
	config.Database = "daemon_app"
	config.Host = " 127.0.0.1"
	config.Password = "#enctypred"
	config.Port = "5432"
	return config
}

// GetConfigurationFromProperties gets database configuration properties from  /config dir
func GetConfigFromProperties() Config {
	file, err := os.Open("./config/config.database_properties.json")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config Config
	if err = decoder.Decode(&config); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	return config
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
