package server

import (
	"encoding/json"
	"net/http"

	"github.com/NenadRadulovic/go-api/storage"
	jwt "github.com/golang-jwt/jwt/v4"
)

type UserLogin struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	AccessToken string `json:"accessToken"`
}

type JWTTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type BuyBookRequest struct {
	BookIds []string `json:"bookIds"`
}

func (s *APIServer) CreateUser(w http.ResponseWriter, r *http.Request) error {
	user, err := s.storage.CreateUser(r.Context(), storage.CreateUserRequest{
		FirstName: r.PostFormValue("firstName"),
		LastName:  r.PostFormValue("lastName"),
	})
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	return nil
}

func (s *APIServer) GetUser(w http.ResponseWriter, r *http.Request) error {
	token, err := s.validateJWT(r.Header.Get("x-lib-token"))
	if err != nil {
		return err
	}

	userId := token.Claims.(jwt.MapClaims)["user_id"].(string)
	user, err := s.storage.GetUserByID(r.Context(), userId)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	return nil
}

func (s *APIServer) BuyBooks(w http.ResponseWriter, r *http.Request) error {
	token, err := s.validateJWT(r.Header.Get("x-lib-token"))
	if err != nil {
		return err
	}

	userId := token.Claims.(jwt.MapClaims)["user_id"].(string)
	user, err := s.storage.GetUserByID(r.Context(), userId)
	if err != nil {
		return err
	}
	var books BuyBookRequest
	_ = json.NewDecoder(r.Body).Decode(&books)

	sum := 0.00
	for _, id := range books.BookIds {
		book, err := s.storage.GetBookById(r.Context(), id)
		if err != nil {
			return err
		}
		sum += book.Price
	}
	remains := (user.Balance - sum)
	if hasMoney := remains < 0; hasMoney {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"purchased": false,
			"message":   "YOU HAVE NO MONEY BITCH",
		})
		return nil
	}

	newUser, err := s.storage.UpdateBalance(r.Context(), userId, remains)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"purchased": true,
		"user":      newUser,
		"message":   "STONKS",
	})
	return nil

}

func (s *APIServer) Login(w http.ResponseWriter, r *http.Request) error {
	var userReq storage.User
	_ = json.NewDecoder(r.Body).Decode(&userReq)
	user, err := s.storage.GetUser(r.Context(), userReq)
	if err != nil {
		return err
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&UserLogin{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		AccessToken: token,
	})
	return nil
}

func (s *APIServer) Logout(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User LOGOUT"))
	return nil
}
