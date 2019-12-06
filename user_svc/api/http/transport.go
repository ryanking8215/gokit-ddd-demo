package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	//"strings"

	//"github.com/go-kit/kit/endpoint"

	"gokit-ddd-demo/user_svc/domain/common"
)

func decodeFindRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idstr := vars["id"]
	if idstr == "" {
		return nil, errors.New("invalid id")
	}

	//idx := strings.Index(r.URL.Path, "users/")
	//if idx < 0 {
	//	return nil, errors.New("invalid url path")
	//}
	//idstr := r.URL.Path[idx+len("users/"):]
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
	//if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
	//	return nil
	//}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, failed error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	if failed == common.ErrNotFound {
		statusCode = http.StatusNotFound
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorWrapper{Error: failed.Error()})
}

type errorWrapper struct {
	Error string `json:"error"`
}
