package handler

import (
	"Apis/internal/users"
	"Apis/pkg/transport"
	"encoding/json"
	"fmt"
	"github.com/JuanFranciscoFI/apiResponse/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

func NewUserHTTPServer(endpoints users.EndPoints) http.Handler {
	r := gin.Default()

	r.POST("/users", transport.GinServer(transport.EndPoint(endpoints.Create), decodePostUserRequest, encodeResponse, encodeError))

	r.GET("/users", transport.GinServer(transport.EndPoint(endpoints.GetAll), decodeGetAllUsersRequest, encodeResponse, encodeError))

	r.GET("/users/:userID", transport.GinServer(transport.EndPoint(endpoints.GetByID), decodeGetUserRequest, encodeResponse, encodeError))

	r.PATCH("/users/:userID", transport.GinServer(transport.EndPoint(endpoints.Update), decodeUpdateUserRequest, encodeResponse, encodeError))

	return r
}

func decodeUpdateUserRequest(c *gin.Context) (interface{}, error) {
	var req users.UpdateRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())

	}

	id, err := strconv.ParseUint(c.Params.ByName("userID"), 10, 64)

	if err != nil {
		return nil, fmt.Errorf("error parsing id: %v", err)
	}

	req.ID = id

	return req, nil
}

func decodePostUserRequest(c *gin.Context) (interface{}, error) {
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	var req users.Request
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return req, nil
}

func decodeGetUserRequest(c *gin.Context) (interface{}, error) {
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	id, err := strconv.ParseUint(c.Params.ByName("userID"), 10, 64)

	if err != nil {
		return nil, err
	}

	return users.GetReq{ID: id}, nil
}

func decodeGetAllUsersRequest(c *gin.Context) (interface{}, error) {
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}

	return nil, nil
}

func tokenVerify(token string) error {
	if os.Getenv("TOKEN") != token {
		return response.Unauthorized("Invalid token")
	}

	return nil
}

func encodeResponse(c *gin.Context, resp interface{}) {
	r := resp.(response.Response)
	c.Header("Content-Type", "application/json")
	c.JSON(r.StatusCode(), resp)
}

func encodeError(c *gin.Context, err error) {
	c.Header("Content-Type", "application/json")
	c.JSON(err.(response.Response).StatusCode(), err)
}
