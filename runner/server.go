package runner

import (
	"os/exec"
)

type ServerInstance struct {
	Running bool
	cmd     *exec.Cmd
}

func (s *ServerInstance) Start(logPath, serverPath string) error {
	cmd := exec.Command("./runner/start.sh >> "+logPath+" 2>&1", serverPath)
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
