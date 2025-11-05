package main

import (
	"auth/auth"
	"auth/middleware"
	"auth/usecase"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	logfilename := fmt.Sprintf("./log/apilog%v", time.Now())

	file, err := os.OpenFile(logfilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error while creating log file")
		return
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := http.NewServeMux()

	mux.HandleFunc("/signup", auth.SignUp)
	mux.HandleFunc("/login", auth.Login)
	mux.HandleFunc("/edit", usecase.EditProfile)

	handler := Chain(
		mux,
		middleware.Loggin,
		middleware.Recovery,
		middleware.Cros,
		middleware.AuthMiddleware,
	)
	log.Println("server started at localhost:8090")
	err = http.ListenAndServe(":8090", handler)
	if err != nil {
		log.Println("error while server start", err)
	}

}

func Chain(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	log.Println("Chain (+)")
	for i := range middleware {
		handler = middleware[i](handler)
	}
	log.Println("Chain (-)")
	return handler
}
