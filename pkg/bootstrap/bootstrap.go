package bootstrap

import (
	"database/sql"
	"log"
)

/*
El bootstrap es el encargado de inicializar las dependencias de la aplicación.
*/

import (
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
}

func NewDB() (*sql.DB, error) {

	dbURL := os.ExpandEnv("$DATABASE_USER:$DATABASE_PASSWORD@tcp($DATABASE_HOST:$DATABASE_PORT)/$DATABASE_NAME")

	db, err := sql.Open("mysql", dbURL)

	if err != nil {
		return nil, err
	}
	return db, nil
}
