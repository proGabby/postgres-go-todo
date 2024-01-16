package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/proGabby/simple_auth_todo_api/pkg/controllers"
	"github.com/proGabby/simple_auth_todo_api/pkg/data/database"
	"github.com/proGabby/simple_auth_todo_api/pkg/middlewares"
	"github.com/proGabby/simple_auth_todo_api/pkg/models"

	"github.com/gorilla/mux"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	connStr, ok := os.LookupEnv("PORT")

	if !ok {
		log.Println("PORT variable not set")
	}
	if connStr == "" {
		log.Fatal("PORT environment variable not set")
	}

	r := mux.NewRouter()

	// Initialize the DB
	db, err := database.InitDB()

	if err != nil {
		log.Fatal(err)
	}

	todoStore := models.NewTodoStore(db)
	userStore := models.NewUserStore(db)

	// Middleware for authentication
	authMiddleware := middlewares.NewAuthMiddleware(*userStore)

	// Middleware for permission
	permissionMiddleware := middlewares.NewPermissionMiddleware(authMiddleware, *todoStore)

	// Initialize controllers
	todoController := controllers.NewTodoController(*todoStore)
	userController := controllers.NewUserController(*userStore)

	// Routes
	r.HandleFunc("/login", userController.LoginUser).Methods("POST")
	r.HandleFunc("/register", userController.RegisterUser).Methods("POST")
	r.HandleFunc("/user/details", authMiddleware.Authenticate(userController.GetUserByToken)).Methods("GET")
	r.HandleFunc("/todos", authMiddleware.Authenticate(todoController.GetTodosByUser)).Methods("GET")
	r.HandleFunc("/todos", authMiddleware.Authenticate(todoController.CreateTodo)).Methods("POST")
	r.HandleFunc("/todos/{id}", authMiddleware.Authenticate(todoController.GetSingleTodo)).Methods("GET")
	r.HandleFunc("/todos/update", authMiddleware.Authenticate(permissionMiddleware.Authorize([]string{"user"}, todoController.UpdateTodo))).Methods("PUT")
	r.HandleFunc("/todos", authMiddleware.Authenticate(permissionMiddleware.Authorize([]string{"admin"}, todoController.DeleteTodo))).Methods("DELETE")
	fmt.Println("before listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
