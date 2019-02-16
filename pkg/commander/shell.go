package commander

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/etherlabsio/errors"
)

// Runnable provides a command than can be run and stopped
type Runnable interface {
	Start() error
	Stop() error
}

// Exec executes a shell command
func Exec(args ...string) error {
	cmd := buildCommand(args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.WithMessagef(err, "exec: failed to pipe RunnableCmd: %v", cmd.Args)
	}
	buf := new(bytes.Buffer)
	err = errors.
		Do(cmd.Start).
		Do(func() error {
			_, err := buf.ReadFrom(stdout)
			return err
		}).
		Do(cmd.Wait).
		Err()
	fmt.Println(buf.String())
	return errors.WithMessagef(err, "exec: failed to execute RunnableCmd: %v", cmd.Args)
}

// RunnableCmd returns a runnable Shell command
type RunnableCmd struct {
	cmd *exec.Cmd
}

func New(args ...string) *RunnableCmd {
	cmd := buildCommand(args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return &RunnableCmd{cmd}
}

// Args returns a concatenated string args
func (s *RunnableCmd) Args() string {
	if s.cmd == nil {
		return ""
	}
	return strings.Join(s.cmd.Args, " ")
}

// Start starts the command execution
func (s *RunnableCmd) Start() error {
	return errors.WithMessage(s.cmd.Start(), "failed to start RunnableCmd: "+s.Args())
}

// Stop terminates the command execution
func (s *RunnableCmd) Stop() error {
	pgid, err := syscall.Getpgid(s.cmd.Process.Pid)
	if err != nil {
		return errors.Wrapf(err, "execute-command: failed to get process group id")
	}

	stop := func(pgid int, signal syscall.Signal) error {
		if err := syscall.Kill(-pgid, signal); err != nil {
			return errors.WithMessage(err, "execute-RunnableCmd: failed to terminate process group id")
		}
		waitErr := s.cmd.Wait()
		err, ok := waitErr.(*exec.ExitError)
		if ok && isSIGTERM(err) {
			return nil
		}
		return errors.WithMessage(waitErr, "waiting for shutdown signal for cmd")
	}

	err = stop(pgid, syscall.SIGTERM)
	return errors.WithMessagef(err, "stop failure: %v", s.cmd.Args)
}

func isSIGTERM(err *exec.ExitError) bool {
	status := err.ProcessState.Sys().(syscall.WaitStatus)
	return status.Signal() == syscall.SIGTERM
}

func buildCommand(args ...string) *exec.Cmd {
	fmt.Println("exec: ", args)
	str := strings.Join(args, " ")
	return exec.Command("sh", "-c", str)
}
