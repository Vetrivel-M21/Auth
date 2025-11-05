package middleware

import (
	"log"
	"net/http"
	"time"
)

func Loggin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%v @ %v @ %v (+)", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("%v @ %v @ %v taken time( %v ) (-)", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}
