package handlers

import (
	"errors"
	"log"
	"net/http"
)

type RecoverHandler struct {
	Action func(error)
	Next   http.Handler
}

func (h *RecoverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		r := recover()
		if r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("Unknown error")
			}

			h.Action(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	h.Next.ServeHTTP(w, r)
}

func Recover(h http.Handler) http.Handler {
	return &RecoverHandler{
		Action: log.Println,
		Next:   h,
	}
}
