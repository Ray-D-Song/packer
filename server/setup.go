package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"ray-d-song.com/packer/dict"
	"ray-d-song.com/packer/server/handler"
	"ray-d-song.com/packer/server/middleware"
	"ray-d-song.com/packer/utils"
)

func init() {
	err := utils.EnsureDirExists(dict.StorageDir)
	if err != nil {
		panic(err)
	}
}

func SetupServer() {
	app := fiber.New(fiber.Config{
		BodyLimit: 200 * 1024 * 1024, // 20MB
	})
	app.Use(middleware.LoggerMiddleware())

	api := app.Group("/api")
	api.Post("/auth/login", handler.UserLogin)
	api.Post("/lib/check", handler.PreCheck)
	api.Post("/lib/publish", handler.Publish)
	api.Post("/lib/download", handler.LibDownload)

	log.Fatal(app.Listen(":7749"))
}
