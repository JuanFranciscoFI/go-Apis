package handler

import (
	"Apis/internal/users"
	"Apis/pkg/transport"
	"context"
	"encoding/json"
	"fmt"
	"github.com/JuanFranciscoFI/apiResponse/response"
	"log"
	"net/http"
	"strconv"
)

func NewUserHTTPServer(ctx context.Context, router *http.ServeMux, endpoints users.EndPoints) {
	router.HandleFunc("/users/", UserServer(ctx, endpoints))
}

func UserServer(ctx context.Context, endpoints users.EndPoints) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		log.Println("URL: ", url)

		path, pathSize := transport.Clean(url)

		params := make(map[string]string)

		if pathSize == 4 && path[2] != "" {
			params["userID"] = path[2]
		}
		t := transport.New(w, r, context.WithValue(ctx, "params", params))

		var end transport.EndPoint
		var deco func(ctx context.Context, r *http.Request) (interface{}, error)

		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3:
				end = transport.EndPoint(endpoints.GetAll)
				deco = decodeGetAllUsersRequest
			case 4:
				end = transport.EndPoint(endpoints.GetByID)
				deco = decodeGetUserRequest
			}
		case http.MethodPost:
			switch pathSize {
			case 3:
				end = transport.EndPoint(endpoints.Create)
				deco = decodePostUserRequest
			}
		case http.MethodPatch:
			switch pathSize {
			case 4:
				end = transport.EndPoint(endpoints.Update)
				deco = decodeUpdateUserRequest
			}
		default:
			InvalidMethod(w)
		}
		if end != nil && deco != nil {
			t.Server(end, deco, encodeResponse, encodeError)
		}
	}

}

func decodeUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req users.UpdateRequest
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&req)

	if err != nil {
		return nil, fmt.Errorf("error decoding request: %v", err)
	}

	params := ctx.Value("params").(map[string]string)

	id, err := strconv.ParseUint(params["userID"], 10, 64)

	if err != nil {
		return nil, fmt.Errorf("error parsing id: %v", err)
	}

	req.ID = id

	return req, nil
}

func decodePostUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req users.Request
	decode := json.NewDecoder(r.Body)
	err := decode.Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	param := ctx.Value("params").(map[string]string)

	id, err := strconv.ParseUint(param["userID"], 10, 64)

	if err != nil {
		return nil, err
	}

	return users.GetReq{ID: id}, nil
}

func decodeGetAllUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "apllication/json")
	resp := err.(response.Response)

	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}

func InvalidMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
