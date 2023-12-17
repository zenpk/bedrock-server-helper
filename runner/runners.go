package runner

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/zenpk/bedrock-server-helper/dal"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Runner struct {
	Db              *dal.Db
	McPath          string
	BaseWorldFolder string
	ServersFolder   string
	BackupsFolder   string
}

// CreateSaveData receives and unzips a world file for later usage
func (r Runner) CreateSaveData(worldId int64, worldFile *multipart.FileHeader, c echo.Context) error {
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	if world.HasSaveData {
		return errors.New("the world already has a save data")
	}
	// TODO transaction
	// world dir
	basePath := r.McPath + "/" + world.Name
	baseWorldPath := basePath + "/" + r.BaseWorldFolder
	serversPath := basePath + "/" + r.ServersFolder
	backupsPath := basePath + "/" + r.BackupsFolder
	if err := runAndOutput(c, "./runner/mkdirs.sh", baseWorldPath, serversPath, backupsPath); err != nil {
		return err
	}
	// copy world zip file
	src, err := worldFile.Open()
	if err != nil {
		return err
	}
	zipFilePath := basePath + r.BaseWorldFolder + worldFile.Filename
	dst, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	unzipDestPath := baseWorldPath + "/" + world.Name + "/"
	if err := runAndOutput(c, "./runner/unzip_rm.sh", zipFilePath, unzipDestPath); err != nil {
		return err
	}
	if err := r.Db.Worlds.SetHasSaveData(worldId, true); err != nil {
		return err
	}
	return endOutput(c)
}

// GetServer downloads a version of server
func (r Runner) GetServer(version string, worldId int64, c echo.Context) error {
	if err := versionNameCheck(version); err != nil {
		return err
	}
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	// TODO transaction
	serverPath := r.McPath + "/" + world.Name + "/" + r.ServersFolder + "/"
	downloadFilePath := serverPath + version + ".zip"
	if err := runAndOutput(c, "./runner/get_server.sh", downloadFilePath, version); err != nil {
		return err
	}
	unzipDestPath := serverPath + version + "/"
	if err := runAndOutput(c, "./runner/unzip_rm.sh", downloadFilePath, unzipDestPath); err != nil {
		return err
	}
	if err := r.Db.Servers.Insert(version, worldId); err != nil {
		return err
	}
	return endOutput(c)
}

// UseServer uses a version of server for a world
func (r Runner) UseServer(serverId, worldId int64, c echo.Context) error {
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	if world.UsingServer == serverId {
		return errors.New("the world is already using the server")
	}
	basePath := r.McPath + "/" + world.Name
	var saveDataPath string
	if world.UsingServer != 0 {
		oldServer, err := r.Db.Servers.SelectById(world.UsingServer)
		if err != nil {
			return err
		}
		saveDataPath = basePath + "/" + r.ServersFolder + "/" + oldServer.Version + "/worlds/" + world.Name
	} else {
		saveDataPath = basePath + "/" + r.BaseWorldFolder + "/" + world.Name
	}
	newServer, err := r.Db.Servers.SelectById(serverId)
	newServerPath := basePath + "/" + r.ServersFolder + "/" + newServer.Version
	// issue: cp -r not behaving as expected, add $5 to fix
	if err := runAndOutput(c, "./runner/use_server.sh", saveDataPath, newServerPath, world.Properties, world.AllowList, world.Name); err != nil {
		return err
	}
	if err := r.Db.Worlds.SetUsingServer(worldId, serverId); err != nil {
		return err
	}
	return endOutput(c)
}

// Backup current world
func (r Runner) Backup(name string, worldId int64, c echo.Context) error {
	var err error
	name, err = r.Db.Backups.ResolveName(name)
	if err != nil {
		return err
	}
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	if world.UsingServer == 0 {
		return errors.New("world is not using a server, there is no need to backup")
	}
	server, err := r.Db.Servers.SelectById(world.UsingServer)
	if err != nil {
		return err
	}
	basePath := r.McPath + "/" + world.Name
	backupPath := basePath + "/" + r.BackupsFolder + "/" + name
	backupPathWithWorldName := backupPath + "/" + world.Name
	saveDataPath := basePath + "/" + r.ServersFolder + "/" + server.Version + "/worlds/" + world.Name
	if err := runAndOutput(c, "./runner/backup.sh", backupPath, backupPathWithWorldName, saveDataPath); err != nil {
		return err
	}
	if err := r.Db.Backups.Insert(name, worldId); err != nil {
		return err
	}
	return endOutput(c)
}

// Restore a backup
func (r Runner) Restore(backupId, worldId int64, ifBackup bool, c echo.Context) error {
	if ifBackup {
		if err := r.Backup("", worldId, c); err != nil {
			return err
		}
	}
	backup, err := r.Db.Backups.SelectById(backupId)
	if err != nil {
		return err
	}
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	server, err := r.Db.Servers.SelectById(world.UsingServer)
	if err != nil {
		return err
	}
	basePath := r.McPath + "/" + world.Name
	backupPath := basePath + "/" + r.BackupsFolder + "/" + backup.Name + "/" + world.Name
	saveDataPath := basePath + "/" + r.ServersFolder + "/" + server.Version + "/worlds/" + world.Name
	if err := runAndOutput(c, "./runner/restore.sh", backupPath, saveDataPath); err != nil {
		return err
	}
	return endOutput(c)
}

// CleanOldBackups deletes backups older than days
func (r Runner) CleanOldBackups(days int64, c echo.Context) error {
	// TODO transaction
	backups, err := r.Db.Backups.SelectDaysBefore(days)
	if err != nil {
		return err
	}
	for _, backup := range backups {
		if err := runAndOutput(c, "./runner/rm_dir.sh", r.BackupsFolder+"/"+backup.Name); err != nil {
			return err
		}
		if err := r.Db.Backups.DeleteById(backup.Id); err != nil {
			return err
		}
	}
	return endOutput(c)
}

func (r Runner) DeleteBackup(worldId, backupId int64, c echo.Context) error {
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	backup, err := r.Db.Backups.SelectById(backupId)
	if err != nil {
		return err
	}
	if err := runAndOutput(c, "./runner/rm_dir.sh", r.McPath+"/"+world.Name+"/"+r.BackupsFolder+"/"+backup.Name); err != nil {
		return err
	}
	if err := r.Db.Backups.DeleteById(backup.Id); err != nil {
		return err
	}
	return endOutput(c)
}

func (r Runner) DeleteServer(worldId, serverId int64, c echo.Context) error {
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	server, err := r.Db.Servers.SelectById(serverId)
	if err != nil {
		return err
	}
	if err := runAndOutput(c, "./runner/rm_dir.sh", r.McPath+"/"+world.Name+"/"+r.ServersFolder+"/"+server.Version); err != nil {
		return err
	}
	if err := r.Db.Servers.DeleteById(server.Id); err != nil {
		return err
	}
	return endOutput(c)
}

// versionNameCheck checks if the version name is "x.x.x.x" format
func versionNameCheck(version string) error {
	numStrs := strings.Split(version, ".")
	err := errors.New("invalid version")
	if len(numStrs) < 4 {
		return err
	}
	for _, numStr := range numStrs {
		num, err2 := strconv.Atoi(numStr)
		if err2 != nil {
			return err
		}
		if num < 0 {
			return err
		}
	}
	return nil
}

// runAndOutput runs a command and outputs the result as server-sent events
func runAndOutput(c echo.Context, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err.Error(), stdErr.String()))
	}
	fmt.Println(stdOut.String())
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	_, err = io.Copy(c.Response(), strings.NewReader(stdOut.String()))
	if err != nil {
		return err
	}
	c.Response().Flush()
	return nil
}

// endOutput ends the server-sent events
func endOutput(c echo.Context) error {
	c.Response().WriteHeader(http.StatusOK)
	_, err := c.Response().Write([]byte("\n\n"))
	return err
}
