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

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3336)/go_course")

	if err != nil {
		return nil, err
	}
	return db, nil
}
