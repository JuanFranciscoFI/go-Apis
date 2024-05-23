package bootstrap

import "log"

/*
El bootstrap es el encargado de inicializar las dependencias de la aplicación.
*/

import (
	"Apis/internal/domain"
	"Apis/internal/users"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)
}

func NewDB() users.DB {
	return users.DB{
		Users: []domain.Users{
			{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "johndoe@gmail.com",
			},
			{
				ID:        2,
				FirstName: "Jane",
				LastName:  "Finn",
				Email:     "janefinn@gmail.com",
			},
			{
				ID:        3,
				FirstName: "Alice",
				LastName:  "Smith",
				Email:     "alicesmith@gmail.com",
			},
		},

		MaxUserID: 3,
	}
}
