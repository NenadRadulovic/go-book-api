package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NenadRadulovic/go-api/storage"
)

func (s *APIServer) CreateBook(w http.ResponseWriter, r *http.Request) error {
	var bookReq storage.CreateBookRequest
	_ = json.NewDecoder(r.Body).Decode(&bookReq)
	log.Println(bookReq)
	book, err := s.storage.CreateBook(r.Context(), bookReq)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
	return nil
}

func (s *APIServer) ListBooks(w http.ResponseWriter, r *http.Request) error {
	books, err := s.storage.ListBooks(r.Context())
	if err != nil {
		return err
	}
	res := map[string][]*storage.Book{
		"books": books,
	}
	json.NewEncoder(w).Encode(res)
	return nil
}

func (s *APIServer) GetBookById(w http.ResponseWriter, r *http.Request) error {
	bookid := s.getUrlParams(r)["id"]
	book, err := s.storage.GetBookById(r.Context(), bookid)
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(book)
	return nil
}

func (s *APIServer) DeleteBook(w http.ResponseWriter, r *http.Request) error {
	bookId := s.getUrlParams(r)["id"]
	res, err := s.storage.DeleteBook(r.Context(), bookId)
	if err != nil {
		return err
	}
	json.NewEncoder(w).Encode(res)
	return nil
}
