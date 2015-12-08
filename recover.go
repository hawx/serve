package serve

import (
	"errors"
	"log"
	"net/http"
)

type recoverHandler struct {
	Action func(error)
	next   http.Handler
}

func (h *recoverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	h.next.ServeHTTP(w, r)
}

func Recover(h http.Handler) http.Handler {
	return &recoverHandler{
		Action: log.Println,
		next:   h,
	}
}
