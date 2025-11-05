package middleware

import (
	"log"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Error MR001", http.StatusInternalServerError)
				log.Printf("Internal Error MR001 %v", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
