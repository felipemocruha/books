package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Category struct {
	ID       int 
	Category string `gorm:"unique"`
}

type Book struct {
	ISBN       string `gorm:"primary_key"`
	Title      string
	Borrowed   bool
	BorrowedBy string
	Categories []Category
}

type DbController interface {
	GetBook(isbn string, book *Book) error
	GetBooks(books *[]Book) error
	CreateBook(isbn, title string) (string, error)
	UpdateBook(isbn, book *Book) (string, error)
	RemoveBook(isbn string) error
	SearcBook(search string, book *[]Book) error
}

type Database struct {
	conn *gorm.DB
}

func (db Database) CreateBook(isbn, title string) (string, error) {
	book := Book{}

	book.ISBN = isbn
	book.Title = title
	book.Borrowed = false
	book.BorrowedBy = ""

	if err = db.conn.Create(&book).Error; err != nil {
		return "", err
	}

	return book.ISBN, nil	
}

func (db Database) GetBook(isbn string, book *Book) error {
	if err := db.conn.Where("isbn = ?", isbn).Find(book).Error; err != nil {
		return err
	}

	return nil
}

func (db Database) GetBooks(books *[]Book) error {
	if err := db.conn.Find(book).Error; err != nil {
		return err
	}

	return nil
}

func (db Database) UpdateBook(isbn, book *Book) (string, error) {

	return isbn, nil
}

func (db Database) RemoveBook(isbn string) error {
	if err := db.conn.Where("isbn = ?", id).Delete(Book{}).Error; err != nil {
		return err
	}

	return nil
}

func (db Database) SearchBook(search string, book *[]Book) error {

	return nil
}

func createDB(conf ServiceConfig) *gorm.DB {
	connStr := "host=%s user=%s dbname=%s sslmode=disable password=%s"
	cmd := fmt.Sprintf(connStr, conf.DbHost, conf.DbUser, conf.DbName, conf.DbPassword)
	db, err := gorm.Open("postgres", cmd)
	if err != nil {
		log.Fatalf("[*] Create database connection error: %v", err)
	}

	db.AutoMigrate(&Book{}, &Category{})
	return db
}
