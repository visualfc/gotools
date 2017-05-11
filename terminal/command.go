// +build !windows

package terminal

import (
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
)

func GetShell() (cmd string, args []string) {
	return "/bin/sh", []string{"-l", "-i"}
}

func Execute(c *exec.Cmd) error {
	f, err := pty.Start(c)
	if err != nil {
		return nil
	}
	go func() {
		for {
			io.Copy(f, os.Stdin)
		}
	}()
	go func() {
		for {
			io.Copy(os.Stdout, f)
		}
	}()
	return c.Wait()
}
