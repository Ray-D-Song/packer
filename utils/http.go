package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code    int16       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseOK(data any) fiber.Map {
	return fiber.Map{
		"code": 200,
		"data": data,
	}
}

func ResponseErr(code int16, message string) fiber.Map {
	return fiber.Map{
		"code":    code,
		"message": message,
	}
}
