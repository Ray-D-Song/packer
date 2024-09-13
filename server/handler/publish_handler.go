package handler

import (
	"io"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"ray-d-song.com/packer/dict"
	"ray-d-song.com/packer/model"
	"ray-d-song.com/packer/utils"
)

var SizeLimit = 1024 * 1024 * 100 * 2 // 200MB

func PreCheck(c *fiber.Ctx) error {
	lib := model.Library{}
	err := c.BodyParser(&lib)
	if err != nil {
		return c.JSON(utils.ResponseErr(500, "Pre-request failed, please try again later"))
	}

	if lib.Size > SizeLimit {
		return c.JSON(utils.ResponseErr(400, "File size exceeds limit"))
	}

	if utils.CheckLibVersionExist(lib.Name, lib.Version) {
		return c.JSON(utils.ResponseErr(400, "Version already exists"))
	}

	return c.JSON(utils.ResponseOK("success"))
}

func Publish(c *fiber.Ctx) error {
	m := c.Queries()
	name := m["name"]
	version := m["version"]
	size := m["size"]

	if name == "" || version == "" || size == "" {
		return c.JSON(utils.ResponseErr(400, "Invalid parameters"))
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(utils.ResponseErr(500, "File upload failed, please try again later"))
	}
	f, err := file.Open()
	if err != nil {
		return c.JSON(utils.ResponseErr(500, "File upload failed, please try again later"))
	}
	defer f.Close()
	dst, err := os.CreateTemp("", "uploaded-*.zip")
	if err != nil {
		return c.JSON(utils.ResponseErr(500, "File upload failed, please try again later"))
	}
	defer os.Remove(dst.Name())

	if _, err := io.Copy(dst, f); err != nil {
		return c.JSON(utils.ResponseErr(500, "File upload failed, please try again later"))
	}

	destDir := path.Join(dict.StorageDir, name, version)
	err = utils.EnsureDirExists(destDir)
	if err != nil {
		return c.JSON(utils.ResponseErr(500, "Create directory failed"))
	}

	if err := utils.Unzip(dst.Name(), destDir); err != nil {
		return c.JSON(utils.ResponseErr(500, "Unable to unzip file"))
	}

	return c.JSON(utils.ResponseOK("success"))
}
