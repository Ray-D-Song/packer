package handler

import (
	"io"
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"ray-d-song.com/packer/dict"
	"ray-d-song.com/packer/model"
	"ray-d-song.com/packer/utils"
)

var SizeLimit = 1024 * 1024 * 100 * 2 // 200MB
var ticketCache []string

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

	// 读取 ~/.packer/perm.key 文件内容
	permFilePath := path.Join(dict.PackerDir, "perm.key")
	permData, err := os.ReadFile(permFilePath)
	if err != nil {
		return c.JSON(utils.ResponseErr(500, "Failed to read perm key file"))
	}
	perm := string(permData)

	// 校验 perm 字段
	if lib.Perm != perm {
		return c.JSON(utils.ResponseErr(400, "Invalid perm key"))
	}

	// 生成新的 UUID
	newUUID := uuid.New().String()
	ticketCache = append(ticketCache, newUUID)

	return c.JSON(utils.ResponseOK(map[string]string{"ticket": newUUID}))
}

func Publish(c *fiber.Ctx) error {
	m := c.Queries()
	name := m["name"]
	version := m["version"]
	size := m["size"]
	p := m["ticket"]

	if name == "" || version == "" || size == "" {
		return c.JSON(utils.ResponseErr(400, "Invalid parameters"))
	}

	found := false
	for i, perm := range ticketCache {
		if perm == p {
			ticketCache = append(ticketCache[:i], ticketCache[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return c.JSON(utils.ResponseErr(400, "Invalid perm key"))
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
