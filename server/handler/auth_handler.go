package handler

import (
	"github.com/gofiber/fiber/v2"
	"ray-d-song.com/packer/model"
)

func UserLogin(c *fiber.Ctx) error {
	form := &model.LoginForm{}
	if err := c.BodyParser(form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if form.Password == "admin" && form.Username == "admin" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Login success",
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid username or password",
	})
}
