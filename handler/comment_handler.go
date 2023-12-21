package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bellaananda/go-postgresql-blog-http.git/models"
	"github.com/bellaananda/go-postgresql-blog-http.git/service"
	"github.com/gorilla/mux"
)

func CreateCommentHandler(commentService service.CommentService, postService service.PostService, userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Extract form values
		user_id := r.Form.Get("user_id")
		post_id := r.Form.Get("post_id")
		content := r.Form.Get("content")

		// Validate and convert form values
		userIDInt, err := strconv.ParseUint(user_id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		postIDInt, err := strconv.ParseUint(post_id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Check if the specified user and post exist
		_, err = userService.GetUserByID(r.Context(), uint(userIDInt))
		if err != nil {
			http.Error(w, "User with the specified ID does not exist", http.StatusBadRequest)
			return
		}
		_, err = postService.GetPostByID(r.Context(), uint(postIDInt))
		if err != nil {
			http.Error(w, "Post with the specified ID does not exist", http.StatusBadRequest)
			return
		}

		// Create a GormComment instance
		comment := models.GormComment{
			UserID:  uint(userIDInt),
			PostID:  uint(postIDInt),
			Content: content,
		}

		// Call the service method to create the comment
		createdComment, err := commentService.CreateComment(r.Context(), comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the created user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(createdComment)
	}
}

func GetAllCommentsHandler(commentService service.CommentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the service method to get the comments
		comments, err := commentService.GetAllComments(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the comments
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

func GetCommentHandler(commentService service.CommentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the comment ID from URL parameters
		idStr := mux.Vars(r)["id"]

		// Convert the ID to int64
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Ensure the ID is not negative
		if id < 0 {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Call the service method to get the comment by id
		comment, err := commentService.GetCommentByID(r.Context(), uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the comment
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comment)
	}
}

func UpdateCommentHandler(commentService service.CommentService, postService service.PostService, userService service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the comment ID from the URL parameters
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Comment ID is missing in URL", http.StatusBadRequest)
			return
		}

		// Parse the ID into an integer
		commentID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		// Extract form values
		userIDStr := r.Form.Get("user_id")
		postIDStr := r.Form.Get("post_id")
		content := r.Form.Get("content")

		// Validate and convert form values
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		postID, err := strconv.ParseUint(postIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// Check if the specified user and post exist
		_, err = userService.GetUserByID(r.Context(), uint(userID))
		if err != nil {
			http.Error(w, "User with specified ID not found", http.StatusBadRequest)
			return
		}
		_, err = postService.GetPostByID(r.Context(), uint(postID))
		if err != nil {
			http.Error(w, "Post with specified ID not found", http.StatusBadRequest)
			return
		}

		// Initialize an empty GormComment
		var updatedComment models.GormComment

		// Set the values for the updated comment
		updatedComment.UserID = uint(userID)
		updatedComment.PostID = uint(postID)
		updatedComment.Content = content

		// Check if the form contains the "content" field
		if content, ok := r.Form["content"]; ok {
			updatedComment.Content = content[0]
		}

		// Set the ID of the comment to be updated
		updatedComment.ID = uint(commentID)

		// Call the service method to update the comment
		existingComment, err := commentService.UpdateCommentByID(r.Context(), uint(commentID), updatedComment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with the updated comment
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingComment)
	}
}

func DeleteCommentHandler(commentService service.CommentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the comment ID from the URL parameters
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			http.Error(w, "Comment ID is missing in URL", http.StatusBadRequest)
			return
		}

		// Parse the ID into an integer
		commentID, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Call the service method to delete the comment
		err = commentService.DeleteCommentByID(r.Context(), uint(commentID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted successfully!"})
	}
}
