package storage

import (
	"context"
	"fmt"
)

type CreateBookRequest struct {
	Title     string  `json:"title"`
	Genre     string  `json:"genre"`
	Author    string  `json:"author"`
	Price     float64 `json:"price"`
	LibraryID int     `json:"library_id"`
}

type Book struct {
	ID      string `json:"book"`
	Title   string
	Genre   string
	Author  string
	Price   float64
	Library Library
}

type Library struct {
	LibraryName string `json:"libraryName"`
	Stars       int    `json:"stars"`
}

func (s *Storage) CreateBook(ctx context.Context, book CreateBookRequest) (*Book, error) {
	row := s.conn.QueryRowContext(ctx, "INSERT INTO books(title, genre, author, price, library_id) VALUES($1, $2, $3, $4, $5) RETURNING id, title, genre, author, price, library_id", book.Title, book.Genre, book.Author, book.Price, book.LibraryID)
	return ScanBook(row)
}

func (s *Storage) ListBooks(ctx context.Context) ([]*Book, error) {
	rows, err := s.conn.QueryContext(ctx, "SELECT books.id, books.title, books.genre, books.author, books.price, library.libraryName, library.stars FROM books LEFT JOIN library ON books.library_id=library.id")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve items %w", err)
	}
	defer rows.Close()
	var books []*Book
	for rows.Next() {
		book, err := ScanBook(rows)
		if err != nil {
			return nil, fmt.Errorf("could not scan book %w", err)
		}
		books = append(books, book)
	}
	return books, nil
}

func (s *Storage) GetBookById(ctx context.Context, bookId string) (*Book, error) {
	row, err := s.conn.QueryContext(ctx, "SELECT * FROM books WHERE ID=$1", bookId)
	if err != nil {
		return nil, fmt.Errorf("could not find by id: %s", bookId)
	}
	defer row.Close()
	row.Next()
	book, err := ScanBook(row)
	if err != nil {
		return nil, fmt.Errorf("error scaning book: %w", err)
	}

	return book, nil
}

func (s *Storage) DeleteBook(ctx context.Context, bookId string) (map[string]bool, error) {

	row, err := s.conn.QueryContext(ctx, "DELETE FROM books WHERE ID=$1 RETURNING ID", bookId)
	if err != nil {
		return map[string]bool{"deleted": false}, fmt.Errorf("could not delete book with id: %s", bookId)
	}
	defer row.Close()

	return map[string]bool{"deleted": true}, nil
}

func ScanBook(s Scanner) (*Book, error) {
	b := &Book{}
	if err := s.Scan(&b.ID, &b.Title, &b.Genre, &b.Author, &b.Price, &b.Library.LibraryName, &b.Library.Stars); err != nil {
		return nil, err
	}
	return b, nil
}
