package http

import (
	"context"
	"encoding/json"
	"gokit-ddd-demo/lib"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func decodeFindRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var userID int64
	var err error
	query := r.URL.Query()
	str := query.Get("user_id")
	if str != "" {
		userID, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0, err
		}
	}
	return userID, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idstr := vars["id"]
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
	Error error `json:"err"`
}
