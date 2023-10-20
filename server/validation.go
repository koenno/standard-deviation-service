package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

var (
	ErrParamNotInteger         = errors.New("parameter must be an integer")
	ErrParamNotPositiveInteger = errors.New("parameter must be a positive integer")
)

//go:generate mockery --name=Handler --srcpkg net/http --case underscore --with-expecter

func validationMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		_, err := paramPositiveInt(r, "requests")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = paramPositiveInt(r, "length")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func paramPositiveInt(r *http.Request, param string) (int, error) {
	requestsStr := r.URL.Query().Get(param)
	value, err := strconv.Atoi(requestsStr)
	if err != nil {
		return 0, fmt.Errorf("%s %w", param, ErrParamNotInteger)
	}
	if value <= 0 {
		return 0, fmt.Errorf("%s %w", param, ErrParamNotPositiveInteger)
	}
	return value, nil
}
