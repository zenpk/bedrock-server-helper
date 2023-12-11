package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zenpk/bedrock-server-helper/dal"
)

type Handlers struct {
	Db *dal.Db
}

func (h Handlers) backupsList(c echo.Context) error {
	backups, err := h.Db.Backups.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, backups)
}

func (h Handlers) versionsList(c echo.Context) error {
	versions, err := h.Db.Versions.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, versions)
}
