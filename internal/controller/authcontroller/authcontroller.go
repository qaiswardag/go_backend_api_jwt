package authcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/qaiswardag/go_backend_api_jwt/database"
	"github.com/qaiswardag/go_backend_api_jwt/internal/logger"
	"github.com/qaiswardag/go_backend_api_jwt/internal/model"
	"github.com/qaiswardag/go_backend_api_jwt/internal/utils"
)

/*
   |--------------------------------------------------------------------------
   | Controller Method Naming Convention
   |--------------------------------------------------------------------------
   |
   | Controller methods: index, create, store, show, edit, update, destroy.
   | Please aim for consistency by using these method names in all controllers.
   |
*/

// Get the user from the context and send it as a response
func Show(w http.ResponseWriter, r *http.Request) {
	fileLogger := logger.FileLogger{}

	user, _ := r.Context().Value("userKey").(model.User)

	response := map[string]interface{}{
		"user": user,
	}

	w.WriteHeader(http.StatusOK)
	fileLogger.LogToFile("AUTH", "Successfully found user and sent response.")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error"})
		fileLogger.LogToFile("AUTH", "Error encoding JSON response")
		return
	}
}

func Destroy(w http.ResponseWriter, r *http.Request) {
	fileLogger := logger.FileLogger{}

	utils.RemoveCookie(w, "session_token", true)
	utils.RemoveCookie(w, "csrf_token", false)

	db, err := database.InitDB()
	if err != nil {
		panic("failed to connect database")
	}

	// Retrieve the user from the context
	sessionUser, okSessionUser := r.Context().Value("sessionUserKey").(model.Session)

	if !okSessionUser {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to retrieve session user from context."})
		fileLogger.LogToFile("AUTH", "Failed to retrieve session user from context.")
		return
	}

	// Delete all other sessions that match the UserID and ServerIP
	if err := db.Exec("DELETE FROM sessions WHERE user_id = ? AND server_ip = ?", sessionUser.UserID, sessionUser.ServerIP).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fileLogger.LogToFile("AUTH", "Failed to delete all other sessions that match the UserID and ServerIP: "+err.Error())
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Internal server error."})

}
