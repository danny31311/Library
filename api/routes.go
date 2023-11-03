package api

import (
	"github.com/gorilla/mux"
	"library/internals/app/handlers"
)

func CreateRoutes(userHandler *handlers.UserHandler, booksHandler *handlers.BooksHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users/create", userHandler.Create).Methods("POST")
	r.HandleFunc("/users/list", userHandler.List).Methods("GET")
	r.HandleFunc("/users/find/{id:[0-9]+}", userHandler.Find).Methods("GET")

	r.HandleFunc("/books/create", booksHandler.Create).Methods("POST")
	r.HandleFunc("/books/list", booksHandler.List).Methods("GET")
	r.HandleFunc("/books/find/{id:[0-9]+}", booksHandler.Find).Methods("GET")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(handlers.NotFound).GetHandler()

	return r

}
