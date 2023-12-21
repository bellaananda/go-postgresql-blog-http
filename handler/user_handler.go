package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
	"github.com/bellaananda/go-postgresql-blog-http.git/service"
	"github.com/gorilla/mux"
)

func CreateUserHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Extract form values
		name := r.Form.Get("name")
		email := r.Form.Get("email")
		password := r.Form.Get("password")
		username := r.Form.Get("username")

		// Create a GormUser instance
		user := models.GormUser{
			Name:     name,
			Email:    email,
			Password: password,
			Username: username,
		}

		// Call the service method to create a user
		createdUser, err := userService.CreateUser(r.Context(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the created user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdUser)
	}
}

func GetAllUsersHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the service method to get the users
		users, err := userService.GetAllUsers(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the users
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUserHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from URL parameters
		idStr := mux.Vars(r)["id"]

		// Convert the ID to int64
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Ensure the ID is not negative
		if id < 0 {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Convert the ID to uint
		uid := uint(id)

		// Call the service method to get the user by id
		user, err := userService.GetUserByID(r.Context(), uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUserHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the URL parameters
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, "User ID is missing in URL", http.StatusBadRequest)
			return
		}

		// Parse the ID into an integer
		userID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Initialize an empty GormUser
		var updatedUser models.GormUser

		// Check if the form contains the "name" field
		if name, ok := r.Form["name"]; ok {
			updatedUser.Name = name[0]
		}

		// Check if the form contains the "email" field
		if email, ok := r.Form["email"]; ok {
			updatedUser.Email = email[0]
		}

		// Check if the form contains the "password" field
		if password, ok := r.Form["password"]; ok {
			updatedUser.Password = password[0]
		}

		// Check if the form contains the "username" field
		if username, ok := r.Form["username"]; ok {
			updatedUser.Username = username[0]
		}

		// Set the ID of the user to be updated
		updatedUser.ID = uint(userID)

		// Call the service method to update the user
		existingUser, err := userService.UpdateUserByID(r.Context(), uint(userID), updatedUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the updated user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingUser)
	}
}

func DeleteUserHandler(userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from the URL parameters
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, "User ID is missing in URL", http.StatusBadRequest)
			return
		}

		// Parse the ID into an integer
		userID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Call the service method to delete the user
		err = userService.DeleteUserByID(r.Context(), uint(userID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
	}
}
