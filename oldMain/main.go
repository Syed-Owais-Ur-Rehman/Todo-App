package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello World")
	app := fiber.New()

	// Defining and getting environment variables
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{} // slice

	// Greeting
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "hello world"})
	})

	// Get all Todos
	app.Get("/api/todos", func(c *fiber.Ctx) error {

		fmt.Println(todos)

		if len(todos) == 0 {
			return c.Status(400).JSON("No Todos added as yet")
		}
		return c.Status(200).JSON(todos)

	})

	// Create a Todo
	app.Post("/api/todo", func(c *fiber.Ctx) error {

		// Remember todo is a pointer to the Todo struct
		todo := &Todo{} // {id:0, completed:false, body:""} by default

		// Handling errors with body parsing
		if err := c.BodyParser(todo); err != nil {
			return err
			// return c.Status(400).JSON(fiber.Map{"error": "error with parsing body"})

		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	//Update a Todo
	app.Put("/api/todo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if strconv.Itoa(todo.ID) == id {
				todos[i].Completed = true

				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not Found"})
	})

	// Delete a Todo
	app.Delete("api/todo/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if strconv.Itoa(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...) // variadic operator (...) will unpack/spread the values.
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not Found"})
	})

	log.Fatal(app.Listen(":" + PORT))

}
