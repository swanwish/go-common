package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

type RouterHandlers interface {
	GetPathPrefix() string
	InitRouter(r *mux.Router)
}

func InitHandlers(handlers []RouterHandlers) {
	r := mux.NewRouter()
	for _, handler := range handlers {
		pathPrefix := handler.GetPathPrefix()
		if pathPrefix != "" {
			s := r.PathPrefix(pathPrefix).Subrouter()
			handler.InitRouter(s)
		} else {
			handler.InitRouter(r)
		}
	}

	// The following register router, so the router will be enabled
	http.Handle("/", r)
}
