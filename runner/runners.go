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
	ServerInstances map[int64]*ServerInstance
	McPath          string
	LogFolder       string
	ServersFolder   string
	BackupsFolder   string
}

// InitDir mkdir for McPath, ServersFolder, BackupsFolder
// Only run one time
func (r Runner) InitDir() error {
	serversPath := r.McPath + "/" + r.ServersFolder
	if err := exec.Command("./runner/mkdir.sh", serversPath).Run(); err != nil {
		return err
	}
	backupsPath := r.McPath + "/" + r.BackupsFolder
	if err := exec.Command("./runner/mkdir.sh", backupsPath).Run(); err != nil {
		return err
	}
	logsPath := r.McPath + "/" + r.LogFolder
	if err := exec.Command("./runner/mkdir.sh", logsPath).Run(); err != nil {
		return err
	}
	return nil
}

// GetServer downloads a version of server
func (r Runner) GetServer(version string, c echo.Context) error {
	if err := versionNameCheck(version); err != nil {
		return err
	}
	// TODO transaction
	serverPath := r.McPath + "/" + r.ServersFolder + "/"
	downloadFilePath := serverPath + version + ".zip"
	if err := runAndOutput(c, "./runner/get_server.sh", downloadFilePath, version); err != nil {
		return err
	}
	unzipDestPath := serverPath + version + "/"
	if err := runAndOutput(c, "./runner/unzip_rm.sh", downloadFilePath, unzipDestPath); err != nil {
		return err
	}
	// create worlds folder
	if err := runAndOutput(c, "./runner/mkdir.sh", unzipDestPath+"worlds"); err != nil {
		return err
	}
	if err := r.Db.Servers.Insert(version); err != nil {
		return err
	}
	return endOutput(c)
}

// UploadSaveData receives and unzips a world file for later usage
func (r Runner) UploadSaveData(worldId int64, worldFile *multipart.FileHeader, c echo.Context) error {
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	if world.HasSaveData {
		return errors.New("the world already has a save data")
	}
	server, err := r.Db.Servers.SelectLatest()
	if err != nil {
		return err
	}
	// TODO transaction
	// copy world zip file
	src, err := worldFile.Open()
	if err != nil {
		return err
	}
	saveDataPath := r.McPath + "/" + r.ServersFolder + "/" + server.Version + "/worlds/"
	zipFilePath := saveDataPath + world.Name + ".zip"
	dst, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	// !IMPORTANT: make sure the zip file contains and only contains one layer of folder
	if err := runAndOutput(c, "./runner/unzip_rm.sh", zipFilePath, saveDataPath); err != nil {
		return err
	}
	// create its backup folder
	worldsBackupPath := r.McPath + "/" + r.BackupsFolder + "/" + world.Name
	if err := runAndOutput(c, "./runner/mkdir.sh", worldsBackupPath); err != nil {
		return err
	}
	if err := r.Db.Worlds.SetHasSaveData(worldId, true); err != nil {
		return err
	}
	if err := r.Db.Worlds.SetUsingServer(worldId, server.Id); err != nil {
		return err
	}
	return endOutput(c)
}

// UseServer uses a version of server for a world
func (r Runner) UseServer(serverId, worldId int64, c echo.Context) error {
	if r.ServerInstances[worldId] != nil && r.ServerInstances[worldId].Running {
		return errors.New("server is running, please stop it first")
	}
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	if world.HasSaveData == false {
		return errors.New("the world does not have a save data")
	}
	if world.UsingServer == serverId {
		return errors.New("the world is already using the server")
	}
	oldServer, err := r.Db.Servers.SelectById(world.UsingServer)
	if err != nil {
		return err
	}
	saveDataPath := r.McPath + "/" + r.ServersFolder + "/" + oldServer.Version + "/worlds/" + world.Name
	newServer, err := r.Db.Servers.SelectById(serverId)
	newSaveDataPath := r.McPath + "/" + r.ServersFolder + "/" + newServer.Version + "/worlds/"
	if err := runAndOutput(c, "./runner/mv.sh", saveDataPath, newSaveDataPath); err != nil {
		return err
	}
	if err := r.Db.Worlds.SetUsingServer(worldId, serverId); err != nil {
		return err
	}
	return endOutput(c)
}

// Backup current world
func (r Runner) Backup(name string, worldId int64, c echo.Context) error {
	if err := r.BackupUtil(name, worldId, c); err != nil {
		return err
	}
	return endOutput(c)
}

// Restore a backup
func (r Runner) Restore(backupId, worldId int64, ifBackup bool, c echo.Context) error {
	if ifBackup {
		if err := r.BackupUtil("", worldId, c); err != nil {
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
	backupPath := r.McPath + "/" + r.BackupsFolder + "/" + world.Name + "/" + backup.Name + "/" + world.Name
	saveDataPath := r.McPath + "/" + r.ServersFolder + "/" + server.Version + "/worlds/" + world.Name
	if err := runAndOutput(c, "./runner/restore.sh", backupPath, saveDataPath); err != nil {
		return err
	}
	return endOutput(c)
}

// CleanOldBackups deletes backups older than days
func (r Runner) CleanOldBackups(worldId, days int64, c echo.Context) error {
	// TODO transaction
	backups, err := r.Db.Backups.SelectDaysBefore(worldId, days)
	if err != nil {
		return err
	}
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	for _, backup := range backups {
		deletePath := r.McPath + "/" + r.BackupsFolder + "/" + world.Name + "/" + backup.Name
		if err := runAndOutput(c, "./runner/rm_dir.sh", deletePath); err != nil {
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
	if err := runAndOutput(c, "./runner/rm_dir.sh", r.McPath+"/"+r.BackupsFolder+"/"+world.Name+"/"+backup.Name); err != nil {
		return err
	}
	if err := r.Db.Backups.DeleteById(backup.Id); err != nil {
		return err
	}
	return endOutput(c)
}

func (r Runner) DeleteServer(serverId int64, c echo.Context) error {
	beingUsed, err := r.Db.Servers.IsInUse(serverId)
	if beingUsed {
		return errors.New("cannot delete the server being used")
	}
	server, err := r.Db.Servers.SelectById(serverId)
	if err != nil {
		return err
	}
	if err := runAndOutput(c, "./runner/rm_dir.sh", r.McPath+"/"+r.ServersFolder+"/"+server.Version); err != nil {
		return err
	}
	if err := r.Db.Servers.DeleteById(server.Id); err != nil {
		return err
	}
	return endOutput(c)
}

func (r Runner) Start(worldId int64) error {
	world, err := r.Db.Worlds.SelectById(worldId)
	if err != nil {
		return err
	}
	if world.UsingServer == 0 {
		return errors.New("no available server")
	}
	server, err := r.Db.Servers.SelectById(world.UsingServer)
	if err != nil {
		return err
	}
	r.ServerInstances[worldId] = &ServerInstance{}
	logPath := r.McPath + "/" + r.LogFolder + "/" + world.Name + ".log"
	serverPath := r.McPath + "/" + r.ServersFolder + "/" + server.Version
	return r.ServerInstances[worldId].Start(logPath, serverPath, world)
}

func (r Runner) Stop(worldId int64) error {
	if r.ServerInstances[worldId] == nil || !r.ServerInstances[worldId].Running {
		return errors.New("server is not running")
	}
	return r.ServerInstances[worldId].Stop()
}

// BackupUtil is used by backup and restore
func (r Runner) BackupUtil(name string, worldId int64, c echo.Context) error {
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
	backupPath := r.McPath + "/" + r.BackupsFolder + "/" + world.Name + "/" + name
	backupPathWithWorldName := backupPath + "/" + world.Name
	saveDataPath := r.McPath + "/" + r.ServersFolder + "/" + server.Version + "/worlds/" + world.Name
	if err := runAndOutput(c, "./runner/backup.sh", backupPath, backupPathWithWorldName, saveDataPath); err != nil {
		return err
	}
	if err := r.Db.Backups.Insert(name, worldId); err != nil {
		return err
	}
	return nil
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
	if c != nil {
		c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
		c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
		_, err = io.Copy(c.Response(), strings.NewReader(stdOut.String()))
		if err != nil {
			return err
		}
		c.Response().Flush()
	}
	return nil
}

// endOutput ends the server-sent events
func endOutput(c echo.Context) error {
	c.Response().WriteHeader(http.StatusOK)
	_, err := c.Response().Write([]byte("\n\n"))
	return err
}
