package router

import (
	"github.com/bellaananda/go-postgresql-blog-http.git/database"
	"github.com/bellaananda/go-postgresql-blog-http.git/handler"
	"github.com/bellaananda/go-postgresql-blog-http.git/repository"
	"github.com/bellaananda/go-postgresql-blog-http.git/service"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	db, errDb := database.RunDatabase()
	if errDb != nil {
		panic("Failed to connect to database")
	}

	router := mux.NewRouter()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, db)

	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository, db)

	commentRepository := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepository, db)

	router.HandleFunc("/api/nicetry", handler.FirstHandler).Methods("GET")
	router.HandleFunc("/api/migrate", handler.MigrateHandler).Methods("GET")
	router.HandleFunc("/api/connect", handler.ConnectHandler).Methods("GET")

	// User routes
	router.HandleFunc("/api/users", handler.CreateUserHandler(*userService)).Methods("POST")                     // create
	router.HandleFunc("/api/users", handler.GetAllUsersHandler(*userService)).Methods("GET")                     // read
	router.HandleFunc("/api/users/{id:[0-9]+}", handler.GetUserHandler(*userService)).Methods("GET")             // read 1
	router.HandleFunc("/api/users/{id:[0-9]+}", handler.UpdateUserHandler(*userService)).Methods("PUT", "PATCH") // update
	router.HandleFunc("/api/users/{id:[0-9]+}", handler.DeleteUserHandler(*userService)).Methods("DELETE")       // delete

	// Post routes
	router.HandleFunc("/api/posts", handler.CreatePostHandler(*postService, *userService)).Methods("POST")                     // create
	router.HandleFunc("/api/posts", handler.GetAllPostsHandler(*postService)).Methods("GET")                                   // read
	router.HandleFunc("/api/posts/{id:[0-9]+}", handler.GetPostHandler(*postService)).Methods("GET")                           // read 1
	router.HandleFunc("/api/posts/{id:[0-9]+}", handler.UpdatePostHandler(*postService, *userService)).Methods("PUT", "PATCH") // update
	router.HandleFunc("/api/posts/{id:[0-9]+}", handler.DeletePostHandler(*postService)).Methods("DELETE")                     // delete

	// Comment routes
	router.HandleFunc("/api/comments", handler.CreateCommentHandler(*commentService, *postService, *userService)).Methods("POST")                     // create
	router.HandleFunc("/api/comments", handler.GetAllCommentsHandler(*commentService)).Methods("GET")                                                 // read
	router.HandleFunc("/api/comments/{id:[0-9]+}", handler.GetCommentHandler(*commentService)).Methods("GET")                                         // read 1
	router.HandleFunc("/api/comments/{id:[0-9]+}", handler.UpdateCommentHandler(*commentService, *postService, *userService)).Methods("PUT", "PATCH") // update
	router.HandleFunc("/api/comments/{id:[0-9]+}", handler.DeleteCommentHandler(*commentService)).Methods("DELETE")                                   // delete

	return router
}
