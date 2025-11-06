package auth

import (
	"auth/common"
	"auth/db"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type SignUpReqData struct {
	Username      string `json:"username" validate:"required"`
	Email         string `json:"email" validate:"required"`
	Password_Hash string `json:"password_hash" validate:"required,min=8"`
}

type LoginReqData struct {
	Identifier string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type LoginUser struct {
	ID           uint
	Email        string
	Username     string
	PasswordHash string
	Role         string
	IsActive     bool
}

func CreateUsersFile(pfilename string, pUser User) error {
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

	// err = writer.Write([]string{pUser.UserName, pUser.Email, pUser.Password_hash})
	// if err != nil {
	// 	log.Printf("Error while writing the file %v", err)
	// 	return err
	// }
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
	var req SignUpReqData
	var lUserData User
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

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		log.Printf("required user information is missing @ ASU004: %v", err)
		http.Error(w, "missing required fields of data", http.StatusBadRequest)
		return
	}

	lUserData.Username = req.Username
	lUserData.Email = req.Email
	lUserData.PasswordHash = req.Password_Hash
	lUserData.Role = "user"
	lUserData.IsActive = true
	lUserData.CreatedAt = time.Now()
	lUserData.UpdatedAt = time.Now()

	//check for existing user
	isExist, lErr := CheckUserExist(lUserData.Username, lUserData.Email)
	if lErr != nil {
		http.Error(w, "Somthing went wrong. Please try again", http.StatusInternalServerError)
		return
	}
	if isExist {
		log.Printf("User already exist with this username/email(ASU005) @ detail : %v", req)
		http.Error(w, "User already exist with this username/email", http.StatusConflict)
		return
	}

	lErr = CreateUser(lUserData)
	if lErr != nil {
		http.Error(w, "Somthing went wrong. Please try again", http.StatusInternalServerError)
		return
	}

	resp.Status = "S"
	resp.Msg = "SignUp Success"

	byteData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error while marshal the data ASU00X :: @ :: %v", err)
		http.Error(w, "Internal Error ASU001", http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(byteData))

}

func CheckUserExist(pusername, pemail string) (bool, error) {
	log.Println(">> CheckUserExist called", "FUNCTION", "CheckUserExist")
	defer log.Println("<< CheckUserExist exited", "FUNCTION", "CheckUserExist")
	var user User

	lErr := db.DB.Table("ST0954_USERS").Where("email = ? and username = ?", pemail, pusername).First(&user).Error
	if lErr != nil {
		if errors.Is(lErr, gorm.ErrRecordNotFound) {
			log.Println("Record already exist(ACUE002)", lErr)
			return false, nil
		} else {
			log.Printf("Error while checking user existance (ACUE001): %v", lErr)
			return false, lErr
		}
	}
	// 2025/11/05 15:06:14 auth.go:125: Error while checking user existance (ACUE001): model value required
	return true, lErr
}

func CreateUser(pUser User) error {
	log.Println(">> CreateUser (+)")
	defer log.Println("<< CreateUser (-)")
	query := `
	INSERT INTO ST0954_USERS 
	(email, password_hash, username, role, is_active, last_login, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	lErr := db.DB.Exec(query,
		pUser.Email,
		pUser.PasswordHash,
		pUser.Username,
		pUser.Role,
		pUser.IsActive,
		pUser.LastLogin,
		pUser.CreatedAt,
		pUser.UpdatedAt,
	).Error

	if lErr != nil {
		log.Printf("Erorr while inserting user data (ACU001) %v", lErr)
		return lErr
	}

	// secondQUERY := ` select id from ST0954_USERS where username = ? and  email = ?`
	var UserId uint

	if lErr = db.DB.Table("ST0954_USERS").Select("id").Where("username = ? and email = ?", pUser.Username, pUser.Email).Scan(&UserId).Error; lErr != nil {
		log.Println("error while find the user id @(ACU002)", lErr)
		return lErr
	}

	var profile UserProfile

	profile.UserID = UserId
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	if lErr = db.DB.Table("ST0954_USERS_PROFILES").Create(&profile).Error; lErr != nil {
		log.Println("error while creating the user profile @(ACU003)", lErr)
		return lErr
	}

	return nil
}

func Login(w http.ResponseWriter, r *http.Request) {

	log.Println(">> Login called", "API", "Login")
	defer log.Println("<< Login exited", "API", "Login")
	if r.Method != http.MethodPost {
		log.Printf("invalid Method %v requires %v for this endpiont %v ALI001", r.Method, http.MethodPost, r.URL.Path)
		http.Error(w, "Invalid Method ", http.StatusMethodNotAllowed)
		return
	}

	var resp LoginResponse
	var req LoginReqData
	var loginUser LoginUser
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error while Reading the request body @ ALI002 @ %v", err)
		http.Error(w, "Faild to read user Input ", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bytes, &req)
	if err != nil {
		log.Printf("Error while Unmarshal the data @ ALI003 @ %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	validator := validator.New()

	if err = validator.Struct(&req); err != nil {
		log.Printf("required user information is missing @ ALI004: %v", err)
		http.Error(w, "missing required fields of data", http.StatusBadRequest)
		return
	}

	query := db.DB.Table("ST0954_USERS").Select("id, email, username, role, is_active")

	if common.IsEmail(req.Identifier) {
		err = query.Where("email = ? and password_hash = ?", req.Identifier, req.Password).First(&loginUser).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("Invalid User credentials (ALI005) @ ", err)
				http.Error(w, "email/password wrong", http.StatusBadRequest)
				return
			} else {
				log.Println("error while User login (ALI006) @ ", err)
				http.Error(w, "Something went wrong. Please try again.", http.StatusBadRequest)
				return
			}
		}
	} else {
		err = query.Where("username = ? and password_hash = ?", req.Identifier, req.Password).First(&loginUser).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("Invalid User credentials (ALI007) @ ", err)
				http.Error(w, "username/password wrong", http.StatusBadRequest)
				return
			} else {
				log.Println("error while User login (ALI008) @ ", err)
				http.Error(w, "Something went wrong. Please try again.", http.StatusBadRequest)
				return
			}
		}
	}

	lLastLogin := time.Now()

	err = db.DB.Table("ST0954_USERS").Where("id = ?", loginUser.ID).Update("last_login", lLastLogin).Error
	if err != nil {
		log.Println("Failed to update the last login (ALI009)")
	}

	token, err := createToken(loginUser.ID, loginUser.Role, loginUser.IsActive)
	if err != nil {
		log.Printf("Error on creating a token for this user (ALI0010) @ %v", loginUser.Username)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	setCookie(token, w)
	resp.Msg = "Login Successfull"
	resp.Status = "S"
	resp.Data.UserName = loginUser.Username
	resp.Data.IsActive = loginUser.IsActive
	resp.Data.Role = loginUser.Role
	resp.Data.UserId = loginUser.ID

	byteData, err := json.Marshal(resp)

	if err != nil {
		log.Println("Error while marshal the data (ALI011)", err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(byteData)
	if writeErr != nil {
		log.Println("Error While Writing Response (ALI012)", writeErr)
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

func createToken(userID uint, role string, isActive bool) (string, error) {
	log.Println(">> createToken called", "Function", "createToken")
	defer log.Println("<< createToken exited", "Function", "createToken")

	lUserID := strconv.FormatUint(uint64(userID), 10)
	claims := jwt.MapClaims{
		"userid":   lUserID,
		"role":     role,
		"isActive": isActive,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}
