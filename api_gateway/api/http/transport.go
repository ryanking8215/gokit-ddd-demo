package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"gokit-ddd-demo/lib"
)

func decodeFindUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vals := r.URL.Query()
	v, ok := vals["with_orders"]
	if ok {
		if len(v) > 0 && v[0] == "true" {
			return true, nil
		}
	}
	return false, nil
}

func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	if idstr == "" {
		return nil, errors.New("invalid id")
	}

	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if err, ok := response.(error); ok {
		errorEncoder(ctx, err, w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, failed error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	err, ok := failed.(lib.Error)
	if ok {
		switch err.Code {
		case lib.ErrNotFound:
			statusCode = http.StatusNotFound
		default:
		}
	} else {
		err = lib.Error{Code: lib.ErrInternal, Message: failed.Error()}
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorWrapper{err})
}

type errorWrapper struct {
	Error error `json:"error"`
}
