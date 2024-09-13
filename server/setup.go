package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"ray-d-song.com/packer/dict"
	"ray-d-song.com/packer/server/handler"
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
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${header:Content-Type} ${path}\n",
	}))
	api := app.Group("/api")
	api.Post("/auth/login", handler.UserLogin)
	api.Post("/lib/check", handler.PreCheck)
	api.Post("/lib/publish", handler.Publish)
	api.Post("/lib/download", handler.LibDownload)
	log.Fatal(app.Listen(":3000"))
}
