package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	//"github.com/go-kit/kit/endpoint"

	"gokit-ddd-demo/order_svc/domain/common"
)

func decodeFindRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var userID int64
	var err error
	query := r.URL.Query()
	str := query.Get("userid")
	if str != "" {
		userID, err = strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0, err
		}
	}
	return userID, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	idx := strings.Index(r.URL.Path, "orders/")
	if idx < 0 {
		return nil, errors.New("invalid url path")
	}
	idstr := r.URL.Path[idx+len("orders/"):]
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
