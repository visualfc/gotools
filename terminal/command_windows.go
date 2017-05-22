// +build windows

package terminal

import (
	"os"
	"os/exec"
)

func checkFiles(names ...string) string {
	for _, name := range names {
		_, err := os.Stat(name)
		if err == nil {
			return name
		}
	}
	return ""
}

func GetShell() (cmd string, args []string) {
	windir := os.Getenv("windir")
	if windir == "" {
		windir = "c:\\windows"
	}
	cmd = checkFiles(windir+"\\Sysnative\\cmd.exe", windir+"\\System32\\cmd.exe")
	return
}

func Execute(c *exec.Cmd) error {
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
