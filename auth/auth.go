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
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is not set")
	}
	TokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("AuthMiddleware: ")
		token, claims, err := jwtauth.FromContext(r.Context())
		println("here 1")
		if err != nil {
			println("here 2")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		if token == nil {
			println("here 3")
			println(token, claims, err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// exp, ok := claims["exp"].(float64)
		// if !ok {
		// 	println("here 4")
		// 	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// 	return
		// }
		// if int64(exp) < time.Now().Unix() {
		// 	println("here 5")
		//
		// 	http.Error(w, "Token has expired", http.StatusUnauthorized)
		// 	return
		// }
		println("here 6")

		next.ServeHTTP(w, r)
	})
}

func RoleMiddleware(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			userRole, ok := claims["role"].(string)
			if !ok || userRole != role {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ExtractUserID(r *http.Request) (string, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return "", err
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return "", fmt.Errorf("user ID claim missing or invalid")
	}

	return userID, nil
}
