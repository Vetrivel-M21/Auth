package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
)

var publicRoute = map[string]bool{
	"/signup": true,
	"/login":  true,
}

func isPublicApi(r *http.Request) bool {
	path := r.URL.Path
	return publicRoute[path]
}

var SecretKey = []byte("your-secret-key")

func validateToken(pJwtToken string) (string, error) {

	token, err := jwt.Parse(pJwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method %v MVT001", token.Header["alg"])
			return nil, fmt.Errorf(" unexpected signing method MVT001 :: @ ::%v : ", token.Header["alg"])
		}
		return SecretKey, nil
	})
	if err != nil {
		log.Printf("Error while parsing the data MVT002  :: @ :: %v", err)
		return "", fmt.Errorf(" Error while parsing the data MVT002 :: @ ::%v : ", err)
	}
	var userIDStr string
	if claim, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		switch v := claim["userid"].(type) {
		case string:
			userIDStr = v
		case float64:
			userIDStr = fmt.Sprintf("%.0f", v) // convert float64 → "123"
		default:
			return "", fmt.Errorf("invalid userid type in token")
		}

		isActive, _ := claim["isActive"].(bool)
		if !isActive {
			return "", fmt.Errorf("inactive user")
		}

		return userIDStr, nil
	}

	return "", fmt.Errorf("invalid Token")

}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if isPublicApi(r) {
			next.ServeHTTP(w, r)
			return
		}

		CokkieData, err := r.Cookie("auth")
		if err != nil {
			http.Error(w, "invalid Session MIV001", http.StatusUnauthorized)
			log.Println("Invalid Session MIV001")
			return
		}

		userid, err := validateToken(CokkieData.Value)
		if err != nil {
			http.Error(w, "invalid Session MIV002", http.StatusUnauthorized)
			log.Println("Invalid Session MIV002")
			return
		}

		r.Header.Set("userid", userid)

		next.ServeHTTP(w, r)

	})
}
