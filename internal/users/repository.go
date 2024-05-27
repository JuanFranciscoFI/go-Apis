package users

import (
	"Apis/internal/domain"
	"database/sql"
	"errors"
	"strings"
)

/*
La capa repository es la encargada de interactuar con la base de datos.
*/

import (
	"context"
	"log"
)

type DB struct {
	Users     []domain.Users
	MaxUserID uint64
}

type (
	Repository interface {
		Create(ctx context.Context, user *domain.Users) error
		GetAll(ctx context.Context) ([]domain.Users, error)
		GetByID(ctx context.Context, id uint64) (*domain.Users, error)
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) error
	}

	repo struct {
		db  *sql.DB
		log *log.Logger
	}
)

func NewRepo(db *sql.DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, u *domain.Users) error {
	sqlQ := "INSERT INTO users (first_name, last_name, email) VALUES (?,?,?)"
	res, err := r.db.Exec(sqlQ, u.FirstName, u.LastName, u.Email)

	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	u.ID = uint64(id)
	r.log.Println("user created with id: ", id)

	return nil
}

func (r *repo) GetAll(ctx context.Context) ([]domain.Users, error) {
	var users []domain.Users

	sqlQ := "SELECT * FROM users"

	rows, err := r.db.Query(sqlQ)

	if err != nil {
		r.log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u domain.Users
		err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
		if err != nil {
			r.log.Println(err.Error())
			return nil, err
		}
		users = append(users, u)
	}

	r.log.Println("users found: ", len(users))
	return users, nil
}

func (r *repo) GetByID(ctx context.Context, id uint64) (*domain.Users, error) {
	var u domain.Users

	sqlQ := "SELECT * FROM users WHERE id = ?"

	err := r.db.QueryRow(sqlQ, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)

	if err != nil {
		r.log.Println(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound{id}
		}
		return nil, err
	}

	r.log.Println("user found with id: ", id)

	return &u, nil
}

func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	var fields []string
	var values []interface{}

	if firstName != nil {
		fields = append(fields, "first_name=?")
		values = append(values, *firstName)
	}

	if lastName != nil {
		fields = append(fields, "last_name=?")
		values = append(values, *lastName)
	}

	if len(fields) == 0 {
		r.log.Println(ErrSliceIsEmpty.Error())
		return ErrSliceIsEmpty
	}

	values = append(values, id)

	sqlQ := "UPDATE users SET " + strings.Join(fields, ", ") + " WHERE id = ?"

	ans, err := r.db.Exec(sqlQ, values...)

	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	row, err := ans.RowsAffected()

	if err != nil {
		r.log.Println(err.Error())
		return err
	}

	if row == 0 {
		return ErrNotFound{ID: id}
	}

	return nil
}
