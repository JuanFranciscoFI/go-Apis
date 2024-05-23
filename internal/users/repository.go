package users

import (
	"Apis/internal/domain"
	"slices"
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
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) (*domain.Users, error)
	}

	repo struct {
		db  DB
		log *log.Logger
	}
)

func NewRepo(db DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, u *domain.Users) error {
	r.db.MaxUserID++
	u.ID = r.db.MaxUserID
	r.db.Users = append(r.db.Users, *u)

	r.log.Println("Repository create")
	return nil
}

func (r *repo) GetAll(ctx context.Context) ([]domain.Users, error) {
	r.log.Println("Repository get all")
	return r.db.Users, nil
}

func (r *repo) GetByID(ctx context.Context, id uint64) (*domain.Users, error) {
	index := slices.IndexFunc(r.db.Users, func(v domain.Users) bool {
		return v.ID == id
	})

	if index == -1 {
		return nil, ErrNotFound{id}
	}

	return &r.db.Users[index], nil
}

func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) (*domain.Users, error) {
	user, err := r.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if firstName != nil {
		user.FirstName = *firstName
	}

	if lastName != nil {
		user.LastName = *lastName
	}

	if email != nil {
		user.Email = *email
	}

	return user, nil
}
