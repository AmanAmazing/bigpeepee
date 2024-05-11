package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	// FIX: Need to have the environment variables loaded here because golang :)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	fmt.Println(os.Getenv("JWT_SECRET_KEY"))
	TokenAuth = jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET_KEY")), nil)
}

func GenerateJWT(claims map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// TODO: Need to add checks for jwt token expiry :)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if claims["role"] != "admin" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ManagerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if claims["role"] != "manager" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request Headers:", r.Header)
		fmt.Println("Request Cookies:", r.Cookies())
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if claims["role"] != "user" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
