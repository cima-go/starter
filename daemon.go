package starter

import (
	"log"
	"os"
	"syscall"

	"github.com/pkg/errors"
	"github.com/sevlyar/go-daemon"
)

var ErrProcessNotExists = errors.New("process not exists")

type daemonWrap struct {
	pid int
	dmx *daemon.Context
	dmp *os.Process
}

func newDaemon(pid, log string) *daemonWrap {
	dmx := &daemon.Context{
		PidFileName: pid,
		LogFileName: log,
	}

	dmp, err := dmx.Search()
	if dmp != nil {
		err = dmp.Signal(syscall.Signal(0))
	}

	run := 0
	if err == nil && dmp != nil {
		run = dmp.Pid
	}

	return &daemonWrap{run, dmx, dmp}
}

func (s *daemonWrap) runs() int {
	return s.pid
}

func (s *daemonWrap) start() (bool, error) {
	if child, err := s.dmx.Reborn(); err != nil {
		return false, err
	} else if child != nil {
		log.Printf("process started, pid = %d", child.Pid)
		return true, nil
	} else {
		return false, nil
	}
}

func (s *daemonWrap) free() {
	if err := s.dmx.Release(); err != nil {
		log.Printf("daemon release failed: %s", err.Error())
	}
}

func (s *daemonWrap) reload() error {
	if s.pid > 0 {
		return s.dmp.Signal(syscall.SIGHUP)
	}

	return ErrProcessNotExists
}

func (s *daemonWrap) shutdown() error {
	if s.pid > 0 {
		return s.dmp.Signal(syscall.SIGTERM)
	}

	return ErrProcessNotExists
}
