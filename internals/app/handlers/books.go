package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"library/internals/app/models"
	"library/internals/app/processors"
	"net/http"
	"strconv"
	"strings"
)

type BooksHandler struct {
	processor *processors.BooksProcessor
}

func NewBooksHandler(processor *processors.BooksProcessor) *BooksHandler {
	handler := new(BooksHandler)
	handler.processor = processor
	return handler
}

func (handler *BooksHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		WrapError(w, err)
		return
	}
	err = handler.processor.CreateBook(newBook)
	if err != nil {
		WrapError(w, err)
		return
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   "",
	}
	WrapOk(w, m)
}

func (handler *BooksHandler) List(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var userIdFilter int64 = 0
	if vars.Get("userid") != "" {
		var err error
		userIdFilter, err = strconv.ParseInt(vars.Get("userid"), 10, 64)
		if err != nil {
			WrapError(w, err)
			return
		}

	}
	list, err := handler.processor.ListBooks(userIdFilter, strings.Trim(vars.Get("book_name"), "\""), strings.Trim(vars.Get("publish_year"), "\""))
	if err != nil {
		WrapError(w, err)
	}
	var m = map[string]interface{}{
		"result": "OK",
		"data":   list,
	}
	WrapOk(w, m)

}

func (handler *BooksHandler) Find(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["id"] == "" {
		WrapError(w, errors.New("missing id"))
		return
	}

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WrapError(w, err)
		return
	}
	book, _ := handler.processor.FindBook(id)
	var m = map[string]interface{}{
		"result": "OK",
		"data":   book,
	}
	WrapOk(w, m)
}
