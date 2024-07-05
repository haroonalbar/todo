package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Text string
	Done bool
}

var (
	todos        []Todo
	loggedInUser string
)

func main() {

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
