package users

import (
	"context"
	"errors"
	"github.com/JuanFranciscoFI/apiResponse/response"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	EndPoints struct {
		Create  Controller
		GetAll  Controller
		GetByID Controller
		Update  Controller
	}

	GetReq struct {
		ID uint64
	}

	Request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	UpdateRequest struct {
		ID        uint64  `json:"id"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
	}
)

func MakeEndPoints(ctx context.Context, s Service) EndPoints {
	return EndPoints{
		Create:  makeCreateEndPoints(s),
		GetAll:  makeGetAllUsersEndPoints(s),
		GetByID: makeGetUserByIDEndPoints(s),
		Update:  makeUpdateUserEndPoints(s),
	}
}

func makeUpdateUserEndPoints(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameIsRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameIsRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email)

		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("updated", nil), nil
	}
}

func makeCreateEndPoints(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)

		if req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameIsRequired.Error())
		}

		if req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameIsRequired.Error())
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("created", user), nil
	}
}

func makeGetUserByIDEndPoints(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)

		user, err := s.GetByID(ctx, req.ID)
		if err != nil {
			if errors.As(err, &ErrNotFound{}) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", user), nil
	}
}

func makeGetAllUsersEndPoints(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("success", users), nil
	}
}
