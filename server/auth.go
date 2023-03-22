package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NenadRadulovic/go-api/storage"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

func (s *APIServer) initAuthRouter(r *mux.Router) {
	router := r.PathPrefix("/auth").Subrouter()
	router.Use(s.authMiddleware)
	router.Methods("GET").Path("/me").Handler(Endpoint{s.GetUser})
	router.Methods("POST").Path("/buy").Handler(Endpoint{s.BuyBooks})
}

func (s *APIServer) generateJWT(user *storage.User) (string, error) {
	secret := "golang-bookstore-api"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"ExpiresAt": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"IssuedAt":  jwt.NewNumericDate(time.Now()),
	})
	tkn, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tkn, nil
}

func (s *APIServer) validateJWT(jwtToken string) (*jwt.Token, error) {
	return jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		// if !t.Valid {
		// 	return nil, fmt.Errorf("invalid token: ")
		// }
		return []byte("golang-bookstore-api"), nil
	})
}

func (s *APIServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-lib-token")
		_, err := s.validateJWT(token)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Not authorized", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
