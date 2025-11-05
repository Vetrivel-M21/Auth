package auth

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var Users []UserInfo

func CreateUsersFile(pfilename string, pUser UserInfo) error {
	log.Println(">> CreateAndUpdateUserFile called", "FUNCTION", "CreateAndUpdateUserFile")
	defer log.Println("<< CreateAndUpdateUserFile exited", "FUNCTION", "CreateAndUpdateUserFile")
	file, err := os.OpenFile(pfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Error while open the file %v", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	err = writer.Write([]string{pUser.UserName, pUser.Email, pUser.Password_hash})
	if err != nil {
		log.Printf("Error while writing the file %v", err)
		return err
	}
	return nil

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	log.Println(">> SignUp called", "API", "SignUp")
	defer log.Println("<< SignUp exited", "API", "SignUp")
	if r.Method != http.MethodPost {
		log.Printf("invalid Method %v requires %v for this endpiont %v ASU001", r.Method, http.MethodPost, r.URL.Path)
		http.Error(w, "Invalid Method ", http.StatusMethodNotAllowed)
		return
	}

	var resp Response
	var req UserInfo
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while Reading the request body @ ASU002 @ %v", err)
		http.Error(w, "Faild to read user Input ", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bytes, &req)
	if err != nil {
		log.Printf("Error while Unmarshal the data @ ASU003 @ %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	if req.UserName == "" || req.Password_hash == "" || req.Email == "" {
		log.Printf("required user information is missing @ ASU004")
		http.Error(w, "missing required fields of data", http.StatusBadRequest)
		return
	}
	err1 := CreateUsersFile("UsersRecord.csv", req)
	if err1 != nil {
		log.Printf("Error in CreateUserFile %v", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}
	Users = append(Users, req)

	resp.Status = "S"
	resp.Msg = "SignUp Success"

	byteData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error while marshal the data ASU00X :: @ :: %v", err)
		http.Error(w, "Internal Error ASU001", http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(byteData))

}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println(">> Login called", "API", "Login")
	defer log.Println("<< Login exited", "API", "Login")
	if r.Method != http.MethodPost {
		log.Printf("invalid Method %v requires %v for this endpiont %v ALI001", r.Method, http.MethodPost, r.URL.Path)
		http.Error(w, "Invalid Method ", http.StatusMethodNotAllowed)
		return
	}

	var resp Response
	var req UserInfo
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while Reading the request body @ ASU002 @ %v", err)
		http.Error(w, "Faild to read user Input ", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bytes, &req)
	if err != nil {
		log.Printf("Error while Unmarshal the data @ ASU003 @ %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	if (req.UserName == "" || req.Email == "") && req.Password_hash == "" {
		log.Printf("required user information is missing @ ASU004")
		http.Error(w, "missing required fields of data", http.StatusBadRequest)
		return
	}
	var isExist bool
	var currUser UserInfo
	for _, user := range Users {
		if (req.Email == user.Email || req.UserName == user.UserName) && req.Password_hash == user.Password_hash {
			isExist = true
			currUser = user
			break
		}
	}

	if !isExist {
		log.Printf("Username/email or password is incorrect @ detail : %v", req)
		http.Error(w, "incorrect username/ password", http.StatusUnauthorized)
		return
	}

	token, err := createToken(currUser.UserName)
	if err != nil {
		log.Printf("Error on creating a token for this user @ %v", currUser)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setCookie(token, w)
	resp.Msg = "Login Successfull"
	resp.Status = "S"

	byteData, err := json.Marshal(resp)

	if err != nil {
		log.Println("Error while marshal the data ", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(byteData)
	if writeErr != nil {
		log.Println("Error While Writing Response", writeErr)
	}
}

var secretKey = []byte("your-secret-key")

func setCookie(token string, w http.ResponseWriter) {
	log.Println(">> setCookie called", "Function", "setCookie")
	defer log.Println("<< setCookie exited", "Function", "setCookie")

	cookie := http.Cookie{
		Name:     "auth",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
	}

	http.SetCookie(w, &cookie)
}

func createToken(username string) (string, error) {
	log.Println(">> createToken called", "Function", "createToken")
	defer log.Println("<< createToken exited", "Function", "createToken")

	claims := jwt.MapClaims{
		"user_name": username,
		"exp":       time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
