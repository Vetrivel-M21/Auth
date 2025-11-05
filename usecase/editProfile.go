package usecase

import (
	"auth/auth"
	"fmt"
	"log"
	"net/http"
)

func EditProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("EditProfile (+)")
	defer log.Println("EditProfile (-)")

	if r.Method != http.MethodPost {
		log.Printf("invalid Method %v requires %v for this endpiont %v ALI001", r.Method, http.MethodPost, r.URL.Path)
		http.Error(w, "Invalid Method ", http.StatusMethodNotAllowed)
		return
	}
	var num1 int
	num2 := 2
	num3 := num2 / num1
	fmt.Println(num3)
	user := auth.Users[0]

	fmt.Fprint(w, user)
}
