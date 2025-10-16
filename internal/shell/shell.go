package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func GetRunningShell() (string, error) {
	if runtime.GOOS == "linux" {
		ppid := os.Getppid()
		cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", ppid)
		
		data, err := os.ReadFile(cmdlineFile)
		if err == nil {
			cmdline := string(data)
			cmdline = strings.TrimPrefix(cmdline, "-")
			parts := strings.Split(cmdline, "\x00")
			if len(parts) > 0 && parts[0] != "" {
				return filepath.Base(parts[0]), nil
			}
		}
	}

	ppid := os.Getppid()
	cmd := exec.Command("ps", "-o", "comm=", "-p", strconv.Itoa(ppid))
	output, err := cmd.Output()
	if err == nil {
		shellName := strings.TrimSpace(string(output))
		shellName = strings.TrimPrefix(shellName, "-")
		shellName = strings.Replace(shellName, "login", "bash", 1)
		if shellName != "" {
			return shellName, nil
		}
	}

	shell := os.Getenv("SHELL")
	if shell != "" {
		return filepath.Base(shell), nil
	}

	return "bash", nil
}

func GetDeskfileExtension() (string, error) {
	shell, err := GetRunningShell()
	if err != nil {
		return ".sh", err
	}

	if shell == "fish" {
		return ".fish", nil
	}
	return ".sh", nil
}
