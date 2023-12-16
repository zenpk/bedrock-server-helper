package main

import (
	"github.com/zenpk/bedrock-server-helper/runner"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zenpk/bedrock-server-helper/dal"
)

type Handlers struct {
	Db     *dal.Db
	Runner *runner.Runner
}

func (h Handlers) worldsList(c echo.Context) error {
	worlds, err := h.Db.Worlds.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, worlds)
}

func (h Handlers) createWorld(c echo.Context) error {
	// TODO name check
	req := struct {
		Name       string `json:"name"`
		Properties string `json:"properties"`
		AllowList  string `json:"allowList"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := h.Db.Worlds.Insert(req.Name, req.Properties, req.AllowList); err != nil {
		return err
	}
	return c.String(http.StatusOK, "ok")
}

func (h Handlers) uploadWorld(c echo.Context) error {
	worldIdStr := c.Param("worldId")
	worldId, err := strconv.ParseInt(worldIdStr, 10, 64)
	if err != nil {
		return err
	}
	file, err := c.FormFile("world")
	if err != nil {
		return err
	}
	return h.Runner.CreateSaveData(worldId, file, c)
}

func (h Handlers) backupsList(c echo.Context) error {
	worldIdStr := c.Param("worldId")
	worldId, err := strconv.ParseInt(worldIdStr, 10, 64)
	if err != nil {
		return err
	}
	backups, err := h.Db.Backups.ListByWorldId(worldId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, backups)
}

func (h Handlers) serversList(c echo.Context) error {
	worldIdStr := c.Param("worldId")
	worldId, err := strconv.ParseInt(worldIdStr, 10, 64)
	if err != nil {
		return err
	}
	versions, err := h.Db.Servers.ListByWorldId(worldId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, versions)
}

func (h Handlers) getServer(c echo.Context) error {
	req := struct {
		WorldId int64  `json:"worldId"`
		Version string `json:"version"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return h.Runner.GetServer(req.Version, req.WorldId, c)
}

func (h Handlers) useServer(c echo.Context) error {
	req := struct {
		ServerId int64 `json:"serverId"`
		WorldId  int64 `json:"worldId"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return h.Runner.UseServer(req.ServerId, req.WorldId, c)
}

func (h Handlers) backup(c echo.Context) error {
	req := struct {
		Name    string `json:"name"`
		WorldId int64  `json:"worldId"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return h.Runner.Backup(req.Name, req.WorldId, c)
}

func (h Handlers) restore(c echo.Context) error {
	req := struct {
		BackupId int64 `json:"backupId"`
		WorldId  int64 `json:"worldId"`
		IfBackup bool  `json:"ifBackup"`
	}{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	return h.Runner.Restore(req.BackupId, req.WorldId, req.IfBackup, c)
}
