package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"library/internals/app/models"
)

type BooksStorage struct {
	dataBasePool *pgxpool.Pool
}

type userBook struct {
	UserId      int64 `db:"userid"`
	Name        string
	Email       string
	Age         int64
	BookId      int64 `db:"bookid"`
	BookName    string
	PublishYear int64
}

func convertJoinedQueryToBook(input userBook) models.Book {
	return models.Book{
		Id:          input.BookId,
		BookName:    input.BookName,
		PublishYear: input.PublishYear,
		Owner: models.User{
			Id:    input.UserId,
			Name:  input.Name,
			Email: input.Email,
			Age:   input.Age,
		},
	}
}

func NewBooksStorage(pool *pgxpool.Pool) *BooksStorage {
	return &BooksStorage{
		dataBasePool: pool,
	}
}

func (storage *BooksStorage) GetBookById(id int64) models.Book {
	query := "SELECT users.id as userid, users.name, users.email, users.age, b.id as bookid, b.book_name, b.publish_year FROM users JOIN books b on users.id = b.user_id WHERE b.id = $1"
	var result []userBook

	err := pgxscan.Select(context.Background(), storage.dataBasePool, &result, query, id)
	if err != nil {
		logrus.Errorln(err)
		return models.Book{}
	}
	if len(result) == 0 {
		return models.Book{}
	}
	return convertJoinedQueryToBook(result[0])

}

func (storage *BooksStorage) GetBooksList(userIdFilter int64, bookNameFilter string, publishYearFilter string) []models.Book {
	query := "SELECT users.id as userid, users.name, users.email, users.age, b.id as bookid, b.book_name, b.publish_year FROM users JOIN books b on users.id = b.user_id WHERE 1=1"
	placeHolderNum := 1
	args := make([]interface{}, 0)
	if userIdFilter != 0 {
		query += fmt.Sprintf(" AND users.ID = $%d", placeHolderNum)
		args = append(args, userIdFilter)
		placeHolderNum++
	}
	if bookNameFilter != "" {
		query += fmt.Sprintf(" AND book_name ILIKE $%d", placeHolderNum)
		args = append(args, fmt.Sprintf("%%%s%%", bookNameFilter))
		placeHolderNum++
	}
	if publishYearFilter != "" {
		query += fmt.Sprintf(" AND publish_year::text LIKE $%d", placeHolderNum)
		args = append(args, fmt.Sprintf("%%%s%%", publishYearFilter))
	}

	var dbResult []userBook
	err := pgxscan.Select(context.Background(), storage.dataBasePool, &dbResult, query, args...)
	if err != nil {
		logrus.Errorln(err)
	}

	result := make([]models.Book, len(dbResult))
	for idx, dbEntity := range dbResult {
		result[idx] = convertJoinedQueryToBook(dbEntity)
	}
	return result
}

func (storage *BooksStorage) CreateBook(book models.Book) error {
	ctx := context.Background()
	tx, err := storage.dataBasePool.Begin(ctx)
	defer func() {
		err = tx.Rollback(context.Background())
		if err != nil {
			logrus.Errorln(err)
		}
	}()

	query := "SELECT id FROM users WHERE ID=$1"
	id := -1

	err = pgxscan.Get(ctx, tx, &id, query, book.Owner.Id)
	if err != nil {
		logrus.Errorln(err)
		err = tx.Rollback(context.Background())
		if err != nil {
			logrus.Errorln(err)
		}
		return err
	}
	if id == -1 {
		return errors.New("user not found")
	}

	insertQuery := "INSERT INTO BOOK(user_id, book_name. publish_year) VALUES ($1, $2, $3)"
	_, err = tx.Exec(context.Background(), insertQuery, book.Owner.Id, book.BookName, book.PublishYear)
	if err != nil {
		logrus.Errorln(err)
		err = tx.Rollback(context.Background())
		if err != nil {
			logrus.Errorln(err)
		}
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		logrus.Errorln(err)
	}
	return err
}
