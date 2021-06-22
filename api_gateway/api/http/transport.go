package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	//"strings"

	//"github.com/go-kit/kit/endpoint"

	"gokit-ddd-demo/user_svc/svc/common"
)

func decodeFindUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vals := r.URL.Query()
	fmt.Println(vals)
	v, ok := vals["with_orders"]
	if ok {
		if len(v) > 0 && v[0] == "true" {
			return true, nil
		}
	}
	return false, nil
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
	if failed == common.ErrNotFound {
		statusCode = http.StatusNotFound
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorWrapper{Error: failed.Error()})
}

type errorWrapper struct {
	Error string `json:"error"`
}
