package runner

import (
	"os/exec"
	"syscall"
)

type ServerInstance struct {
	Running bool
	cmd     *exec.Cmd
}

func (s *ServerInstance) Start(logPath, serverPath string) error {
	cmd := exec.Command("./runner/start.sh", serverPath, logPath)
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
	if err := s.cmd.Process.Kill(); err != nil {
		return err
	}
	s.cmd = nil
	s.Running = false
	return nil
}
