package main

import (
	"bufio"
	"errors"
	"github.com/zenpk/bedrock-server-helper/cron"
	"github.com/zenpk/bedrock-server-helper/runner"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zenpk/bedrock-server-helper/dal"
)

type Handlers struct {
	Db     *dal.Db
	Runner *runner.Runner
	Cron   *cron.Cron
}

func (h Handlers) serversList(c echo.Context) error {
	versions, err := h.Db.Servers.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, versions)
}

func (h Handlers) getServer(c echo.Context) error {
	req := &struct {
		Version string `json:"version"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	return h.Runner.GetServer(req.Version, c)
}

func (h Handlers) useServer(c echo.Context) error {
	req := &struct {
		ServerId int64 `json:"serverId"`
		WorldId  int64 `json:"worldId"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	return h.Runner.UseServer(req.ServerId, req.WorldId, c)
}

func (h Handlers) worldsList(c echo.Context) error {
	worlds, err := h.Db.Worlds.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, worlds)
}

func (h Handlers) createWorld(c echo.Context) error {
	req := &struct {
		Name       string `json:"name"`
		Properties string `json:"properties"`
		AllowList  string `json:"allowList"`
		ServerId   int64  `json:"serverId"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	if strings.Contains(req.Name, " ") || strings.Contains(req.Name, "/") || strings.Contains(req.Name, "\\") || strings.Contains(req.Name, ".") || strings.Contains(req.Name, ":") || strings.Contains(req.Name, "*") || strings.Contains(req.Name, "?") || strings.Contains(req.Name, "\"") || strings.Contains(req.Name, "<") || strings.Contains(req.Name, ">") || strings.Contains(req.Name, "|") {
		return errors.New("cannot use this world name")
	}
	if req.ServerId <= 0 {
		return errors.New("invalid server id")
	}
	if err := h.Db.Worlds.Insert(req.Name, req.Properties, req.AllowList, req.ServerId); err != nil {
		return err
	}
	return nil
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
	return h.Runner.UploadSaveData(worldId, file, c)
}

func (h Handlers) deleteServer(c echo.Context) error {
	req := &struct {
		ServerId int64 `json:"serverId"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	return h.Runner.DeleteServer(req.ServerId, c)
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

func (h Handlers) backup(c echo.Context) error {
	req := &struct {
		Name    string `json:"name"`
		WorldId int64  `json:"worldId"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	return h.Runner.Backup(req.Name, req.WorldId, c)
}

func (h Handlers) restore(c echo.Context) error {
	req := &struct {
		BackupId int64 `json:"backupId"`
		WorldId  int64 `json:"worldId"`
		IfBackup bool  `json:"ifBackup"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	return h.Runner.Restore(req.BackupId, req.WorldId, req.IfBackup, c)
}

func (h Handlers) deleteBackup(c echo.Context) error {
	req := &struct {
		WorldId  int64 `json:"worldId"`
		BackupId int64 `json:"backupId"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	return h.Runner.DeleteBackup(req.WorldId, req.BackupId, c)
}

func (h Handlers) start(c echo.Context) error {
	req := &struct {
		WorldId int64 `json:"worldId"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	return h.Runner.Start(req.WorldId)
}

func (h Handlers) stop(c echo.Context) error {
	req := &struct {
		WorldId int64 `json:"worldId"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	return h.Runner.Stop(req.WorldId)
}

func (h Handlers) isRunning(c echo.Context) error {
	worldIdStr := c.Param("worldId")
	worldId, err := strconv.ParseInt(worldIdStr, 10, 64)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, h.Runner.ServerInstances[worldId] != nil)
}

func (h Handlers) getLog(c echo.Context) error {
	const maxLine = 1000
	worldIdStr := c.Param("worldId")
	worldId, err := strconv.ParseInt(worldIdStr, 10, 64)
	if err != nil {
		return err
	}
	startLineStr := c.Param("startLine")
	startLine, err := strconv.ParseInt(startLineStr, 10, 64)
	if err != nil {
		return err
	}
	world, err := h.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	logPath := h.Runner.McPath + "/" + h.Runner.LogFolder + "/" + world.Name + ".log"
	file, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer file.Close()
	var readLines []string
	scanner := bufio.NewScanner(file)
	currentLine := int64(1)
	for scanner.Scan() {
		if currentLine >= startLine {
			readLines = append(readLines, scanner.Text())
		}
		currentLine++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(readLines) > maxLine {
		readLines = readLines[len(readLines)-maxLine:]
	}
	return c.JSON(http.StatusOK, readLines)
}

func (h Handlers) cronList(c echo.Context) error {
	worldIdStr := c.Param("worldId")
	worldId, err := strconv.ParseInt(worldIdStr, 10, 64)
	if err != nil {
		return err
	}
	crons, err := h.Db.Crons.SelectByWorldId(worldId)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, crons)
}

func (h Handlers) createCron(c echo.Context) error {
	req := &struct {
		JobName    string `json:"jobName"`
		WorldId    int64  `json:"worldId"`
		Parameters string `json:"parameters"`
		Cron       string `json:"cron"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := h.Db.Crons.Insert(req.JobName, req.WorldId, req.Parameters, req.Cron); err != nil {
		return err
	}
	return h.Cron.RefreshCron()
}

func (h Handlers) deleteCron(c echo.Context) error {
	req := &struct {
		Id int64 `json:"id"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := h.Db.Crons.DeleteById(req.Id); err != nil {
		return err
	}
	return h.Cron.RefreshCron()
}
