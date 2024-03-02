package main

import (
	"fly-migrations/views"
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	conn := createConnection()

	app := fiber.New(fiber.Config{})
	// app.Use(logger.New())
	app.Static("/public", "./public")

	app.Get("/", func(c *fiber.Ctx) error {
		root := views.RootLayout("Skill Issue", GetJokes(conn))

		return render(root, c)
	})
	app.Post("/create", func(c *fiber.Ctx) error {
		var body jokeFormPayload
		if err := c.BodyParser(&body); err != nil {
			_form := views.JokeForm("", err.Error())
			return render(_form, c)
		}

		if body.Joke == "" {
			_form := views.JokeForm("", "Missing joke")
			return render(_form, c)
		}

		if _, err := conn.Exec("insert into jokes (joke) values ($1)", body.Joke); err != nil {
			log.Println("error creating joke", err)

			_form := views.JokeForm(body.Joke, "Error creating joke")
			return render(_form, c)
		}

		_form := views.JokeForm("", "")
		_list := views.JokeList(GetJokes(conn))

		c.Response().Header.SetContentType("text/html")
		_form.Render(c.Context(), c.Response().BodyWriter())
		_list.Render(c.Context(), c.Response().BodyWriter())

		return nil
	})

	app.Listen(":8080")
}

type jokeFormPayload struct {
	Joke string `form:"joke"`
}

func render(component templ.Component, c *fiber.Ctx) error {
	c.Response().Header.SetContentType("text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func createConnection() *sqlx.DB {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panicf("failed to open connection: %v\n", err)
	}

	return db
}

func GetJokes(conn *sqlx.DB) []string {
	var jokes []string
	err := conn.Select(&jokes, "select joke from jokes order by id desc")

	if err != nil {
		log.Println("error fetching jokes", err)
		return []string{}
	}

	return jokes
}
