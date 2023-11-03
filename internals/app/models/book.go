package models

type Book struct {
	Id          int64  `json:"id"`
	BookName    string `json:"book_name"`
	PublishYear int64  `json:"publish_year"`
	Owner       User   `json:"owner"`
}
