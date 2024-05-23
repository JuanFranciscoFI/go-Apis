package users

import (
	"Apis/internal/domain"
	"context"
	"log"
)

/*
Esta capa es la encargada de manejar la lógica de negocio de la aplicación.
En este caso, se encarga de crear un nuevo usuario y de obtener todos los usuarios.
*/

type (
	Service interface {
		Create(ctx context.Context, firstName, lastName, email string) (*domain.Users, error)
		GetAll(ctx context.Context) ([]domain.Users, error)
		GetByID(ctx context.Context, id uint64) (*domain.Users, error)
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) (*domain.Users, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, firstName, lastName, email string) (*domain.Users, error) {
	u := domain.Users{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	err := s.repo.Create(ctx, &u)
	if err != nil {
		return nil, err
	}

	s.log.Println("Service create")
	return &u, nil
}

func (s *service) GetByID(ctx context.Context, id uint64) (*domain.Users, error) {
	s.log.Println("Service get by id")
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context) ([]domain.Users, error) {
	s.log.Println("Service get all")
	return s.repo.GetAll(ctx)
}

func (s *service) Update(ctx context.Context, id uint64, firstName, lastName, email *string) (*domain.Users, error) {
	s.log.Println("Service update")
	return s.repo.Update(ctx, id, firstName, lastName, email)
}
