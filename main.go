package main

import (
	"fmt"
	"os"

	"github.com/abe27/oracle/api/configs"
	"github.com/abe27/oracle/api/controllers"
	"github.com/abe27/oracle/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	configs.USERNAME = os.Getenv("DB_USERNAME")
	configs.PASSWORD = os.Getenv("DB_PASSWORD")
	configs.HOST = os.Getenv("DB_HOST")
	configs.DATABASE = os.Getenv("DB_DATABASE")
	configs.REST_URL = os.Getenv("REST_URL")
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusCreated).JSON("Hello, world!")
	})

	app.Post("/carton", func(c *fiber.Ctx) error {
		var obj models.CartonForm
		err := c.BodyParser(&obj)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON("error")
		}

		go controllers.FetchData(&obj)
		return c.Status(fiber.StatusCreated).JSON(&obj.SerialNo)
	})

	app.Get("/carton/search", func(c *fiber.Ctx) error {
		serial_no := c.Query("serial_no")
		if serial_no == "" {
			return c.Status(fiber.StatusBadRequest).JSON("Not Allow!")
		}
		isFound := controllers.FetchDataBySerialNo(serial_no)
		fmt.Printf("serial_no: %s is: %v\n", serial_no, isFound)
		if isFound {
			return c.Status(fiber.StatusFound).JSON(serial_no)
		}
		return c.Status(fiber.StatusNotFound).JSON(serial_no)
	})

	/// router for stock
	app.Get("/stock", controllers.FetchAllStock)
	app.Listen(":4000")
}
