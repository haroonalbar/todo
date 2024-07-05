package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Todo struct {
	Text string
	Done bool
}

var (
	todos        []Todo
	loggedInUser string
)

// key used for signing and verifying JWT tokens
var secretKey = []byte("")

func main() {
	// Load environment variables from .env file if it exists
	// This is helpful for avoiding hardcoding sensitive information like API keys
	godotenv.Load()
	secretKeyEnv := os.Getenv("SECRET_KEY")

	// Add the secret key to the global secretKey variable
	secretKey = []byte(secretKeyEnv)

	// Create a new router with default middleware
	router := gin.Default()
	// Serve static files from the "static" directory
	router.Static("/static", "./static")
	// Load HTML templates from the "templates" directory
	router.LoadHTMLGlob("templates/*")

	// This route handles the root request ("/")
	// It serves the "index.html" template with the following data:
	// - "Todos": the list of todos
	// - "LoggedIn": a boolean indicating whether a user is logged in or not
	// - "Username": the username of the logged in user
	router.GET("/", func(ctx *gin.Context) {
		// Serve the "index.html" template with the list of todos and login status
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"Todos":    todos,
			"LoggedIn": loggedInUser != "",
			"Username": loggedInUser,
		})
	})

	// Adds a new todo item to the list of todos, and redirects to the home page
	router.POST("/add", func(ctx *gin.Context) {
		// Get the text of the todo item from the form data
		text := ctx.PostForm("todo")
		// Create a new todo item with the given text and set it as not done
		todo := Todo{Text: text, Done: false}
		// Append the new todo item to the list of todos
		todos = append(todos, todo)
		// Redirect the user to the home page
		ctx.Redirect(http.StatusSeeOther, "/")
	})

	// This route handles the POST request to toggle the "Done" status of a todo item.
	// It expects the "index" form parameter to specify the index of the todo item to toggle.
	//
	// Parameters:
	// - "index": a string representing the index of the todo item to toggle.
	router.POST("/toggle", func(ctx *gin.Context) {
		// Get the index of the todo item to toggle from the form data
		index := ctx.PostForm("index")

		// Toggle the "Done" status of the todo item at the specified index
		toggleIndex(index)

		// Redirect the user to the home page
		ctx.Redirect(http.StatusSeeOther, "/")
	})

	// Start the Gin router on port 8080
	// This blocks until the server is shutdown, so it should be the last line of the main function
	//
	// Parameters:
	// - ":8080": the port to listen on
	router.Run(":8080")
}

// toggleIndex toggles the "Done" status of a todo item in the "todos" list
// based on the index provided in the "s" parameter.
//
// Parameters:
// - s: a string representing the index of the todo item to toggle.
func toggleIndex(s string) {
	// Convert the string index to an integer.
	i, _ := strconv.Atoi(s)

	// If the index is within the valid range of the "todos" list, toggle the
	// "Done" status of the corresponding todo item.
	if i >= 0 && i < len(todos) {
		todos[i].Done = !todos[i].Done
	}
}

func createToken(usename string) (string, error) {
	// Create a new JWT token with the specified claims
	claims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": usename,                              // subject (user identifier)
		"iss": "todo-app",                           // issuer
		"aud": getRole(usename),                     // audience (user role)
		"exp": time.Now().Add(2 * time.Hour).Unix(), // expiration time
		"iat": time.Now().Unix(),                    // issued at
	})

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)

	// return tokenString,nil
	panic("not used")
}

func getRole(usename string) string {
	return ""
}
