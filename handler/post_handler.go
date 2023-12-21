package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
	"github.com/bellaananda/go-postgresql-blog-http.git/service"
	"github.com/gorilla/mux"
)

func CreatePostHandler(postService service.PostService, userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Extract form values
		user_id := r.Form.Get("user_id")
		title := r.Form.Get("title")
		content := r.Form.Get("content")

		// Convert user_id to uint
		user_id_int, err := strconv.ParseUint(user_id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Check if the specified user ID exists
		_, err = userService.GetUserByID(r.Context(), uint(user_id_int))
		if err != nil {
			http.Error(w, "User with the specified ID does not exist", http.StatusBadRequest)
			return
		}

		// Create a GormPost instance
		post := models.GormPost{
			UserID:  uint(user_id_int),
			Title:   title,
			Content: content,
		}

		// Call the service method to create a post
		createdPost, err := postService.CreatePost(r.Context(), post)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the created post
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdPost)
	}
}

func GetAllPostsHandler(postService service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the service method to get the posts
		posts, err := postService.GetAllPosts(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the posts
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func GetPostHandler(postService service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the post ID from URL parameters
		idStr := mux.Vars(r)["id"]

		// Convert the ID to int64
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Ensure the ID is not negative
		if id < 0 {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Call the service method to get the post
		post, err := postService.GetPostByID(r.Context(), uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

func UpdatePostHandler(postService service.PostService, userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the post ID from the URL parameters
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Post ID is missing in URL", http.StatusBadRequest)
			return
		}

		// Parse the ID into an integer
		postID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Extract form values
		userIDStr := r.Form.Get("user_id")
		title := r.Form.Get("title")
		content := r.Form.Get("content")

		// Convert user ID to uint
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Check if the specified user ID exists
		_, err = userService.GetUserByID(r.Context(), uint(userID))
		if err != nil {
			http.Error(w, "User with specified ID not found", http.StatusBadRequest)
			return
		}

		// Initialize an empty GormPost
		var updatedPost models.GormPost

		// Set the values for the updated post
		updatedPost.UserID = uint(userID)
		updatedPost.Title = title
		updatedPost.Content = content

		// Check if the form contains the "title" field
		if title, ok := r.Form["title"]; ok {
			updatedPost.Title = title[0]
		}

		// Check if the form contains the "content" field
		if content, ok := r.Form["content"]; ok {
			updatedPost.Content = content[0]
		}

		// Set the ID of the post to be updated
		updatedPost.ID = uint(postID)

		// Call the service method to update the post
		existingPost, err := postService.UpdatePostByID(r.Context(), uint(postID), updatedPost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the updated post
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingPost)
	}
}

func DeletePostHandler(postService service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the post ID from the URL parameters
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, "Post ID is missing in URL", http.StatusBadRequest)
			return
		}

		// Parse the ID into an integer
		postID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Call the service method to delete the post
		err = postService.DeletePostByID(r.Context(), uint(postID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully!"})
	}
}
