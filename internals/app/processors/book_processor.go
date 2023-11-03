package processors

import (
	"errors"
	"library/internals/app/db"
	"library/internals/app/models"
)

type BooksProcessor struct {
	storage *db.BooksStorage
}

func NewBooksProcessor(storage *db.BooksStorage) *BooksProcessor {
	processor := new(BooksProcessor)
	processor.storage = storage
	return processor
}

func (processor *BooksProcessor) CreateBook(book models.Book) error {
	if book.BookName == "" {
		return errors.New("book name is empty")
	}
	if book.PublishYear <= 1500 {
		return errors.New("publish year must be greater than 1500")
	}
	if book.Owner.Id <= 0 {
		return errors.New("owner id shall be filled")
	}
	return processor.storage.CreateBook(book)

}

func (processor *BooksProcessor) FindBook(id int64) (models.Book, error) {
	book := processor.storage.GetBookById(id)
	if book.Id != id {
		return book, errors.New("book not found")
	}
	return book, nil
}

func (processor *BooksProcessor) ListBooks(userId int64, bookNameFilter string, publishYearFilter string) ([]models.Book, error) {
	return processor.storage.GetBooksList(userId, bookNameFilter, publishYearFilter), nil

}
