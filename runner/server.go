package runner

import (
	"github.com/zenpk/bedrock-server-helper/dal"
	"os/exec"
	"syscall"
)

type ServerInstance struct {
	Running bool
	cmd     *exec.Cmd
}

func (s *ServerInstance) Start(logPath, serverPath string, world dal.Worlds) error {
	cmd := exec.Command("./runner/start.sh", serverPath, logPath, world.Properties, world.AllowList)
	// Make the command a leader of a new process group
	// This will allow us to kill all related processes in this process group later
	// Linux specific
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	s.cmd = cmd
	s.Running = true
	return nil
}

func (s *ServerInstance) Stop() error {
	// Linux specific
	pgid, err := syscall.Getpgid(s.cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGTERM)
	}
	s.cmd = nil
	s.Running = false
	return nil
}
