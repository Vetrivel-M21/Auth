package usecase

import (
	"auth/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func UpdateProfile(userId uint, payload map[string]interface{}) error {
	var allowFeild = map[string]bool{
		"first_name": true, "last_name": true, "phone": true, "dob": true, "gender": true, "address": true, "city": true, "country": true, "profile_image_url": true,
		"bio": true,
	}

	saveData := make(map[string]interface{})

	for key, val := range payload {
		if allowFeild[key] {
			saveData[key] = val
		}
	}

	if len(saveData) == 0 {
		return fmt.Errorf(" No valid feild to update. ")
	}
	// log.Printf("Userid %v  update their profile %v", userId, saveData)

	tx := db.DB.Table("ST0954_USERS_PROFILES").Where("user_id = ?", userId).Updates(saveData)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return fmt.Errorf("no profile found for user_id: %v", userId)
	}
	return nil
}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	log.Println("EditProfile (+)")
	defer log.Println("EditProfile (-)")

	if r.Method != http.MethodPatch {
		log.Printf("invalid Method %v requires %v for this endpiont %v UEP001", r.Method, http.MethodPost, r.URL.Path)
		http.Error(w, "Invalid Method ", http.StatusMethodNotAllowed)
		return
	}
	lUserID := r.Header.Get("userid")
	log.Println("User Edit Profile :", lUserID)
	u64, err := strconv.ParseUint(lUserID, 10, 64) // base 10, 64-bit unsigned
	if err != nil {
		log.Println("Invalid ID:", err)
		return
	}

	var payload = make(map[string]interface{})

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println("Error while reading the request body (UEP002) @ ", err)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	err = UpdateProfile(uint(u64), payload)
	if err != nil {
		log.Println("Error on Update profile @ (UEP003)", err)
		http.Error(w, "Somthing Went wrong (UEP003)", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Update succesful")

}
