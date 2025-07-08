package login

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adohong4/carZone/core"
	"github.com/adohong4/carZone/helpers"
	"github.com/adohong4/carZone/models"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		// http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		core.SendErrorResponse(w, core.NewBadRequestError("Invalid request body").ErrorResponse)
		return
	}

	valid := (credentials.UserName == "admin" && credentials.Password == "admin123")
	if !valid {
		//http.Error(w, "Incorrect Username or Password", http.StatusUnauthorized)
		core.SendErrorResponse(w, core.NewAuthFailureError("Incorrect Username or Password").ErrorResponse)
		return
	}

	tokenString, err := helpers.GenerateToken(credentials.UserName)
	if err != nil {
		log.Println("Error Generating token: ", err)
		//http.Error(w, "Failed to Username or Password", http.StatusUnauthorized)
		core.SendErrorResponse(w, core.NewAuthFailureError("Failed to Username or Password").ErrorResponse)
		return
	}

	response := map[string]string{"token": tokenString}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
	core.NewOK("Login successful", response).Send(w)

}
