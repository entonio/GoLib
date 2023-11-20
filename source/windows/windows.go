package windows

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"golib/arrays"
	"golib/charsets"
	"golib/fsx"
	"golib/log"
)

func createCommand(dir string, printOutput bool, hide []string, cmdName string, cmdArgs ...string) *exec.Cmd {
	command := exec.Command(cmdName, cmdArgs...)
	if len(dir) > 0 {
		command.Dir = dir
		if printOutput {
			log.Info("CD %s", dir)
		}
	}
	if printOutput {
		log.Info("%s %s", cmdName, strings.Join(arrays.Replace(cmdArgs, hide, "*"), " "))
	}
	return command
}

func Exec(dir string, printOutput bool, hide []string, acceptCodes []int, cmdName string, cmdArgs ...string) (output string, err error) {
	output, _, err = ExecL(dir, "", printOutput, hide, acceptCodes, cmdName, cmdArgs...)
	return
}

func ExecL(dir string, logFile string, printOutput bool, hide []string, acceptCodes []int, cmdName string, cmdArgs ...string) (output string, logs string, err error) {
	var errs []error
	if len(logFile) > 0 {
		os.Remove(logFile)
	}
	command := createCommand(dir, printOutput, hide, cmdName, cmdArgs...)
	stdout, err := command.CombinedOutput()
	if err != nil {
		accepted := arrays.AnyMatch(acceptCodes, func(code int) bool {
			return err.Error() == fmt.Sprintf("exit status %d", code)
		})
		if !accepted {
			errs = append(errs, err)
		}
	}
	output, err = charsets.Decode(stdout, charmap.Windows1252)
	if err != nil {
		errs = append(errs, err)
	}
	if len(logFile) > 0 {
		logsBytes, _ := fsx.ReadBytes(logFile)
		if logsBytes != nil {
			logs = strings.TrimSuffix(string(logsBytes), "\n")
		}
	}
	errorMessage := arrays.Print1d(errs, ", ")
	if printOutput {
		if len(output) > 0 {
			log.Info(output)
		} else if len(errs) > 0 {
			log.Error("%s error(s): %s", cmdName, errorMessage)
		} else if len(logs) > 0 {
			log.Info(logs)
		} else {
			log.Debug("%s produced no output or logs", cmdName)
		}
	}
	if len(errs) > 0 {
		err = errors.New(errorMessage)
	} else {
		err = nil
		log.Debug("%s completed successfully", cmdName)
	}
	return
}

func Start(dir string, printOutput bool, hide []string, cmdName string, cmdArgs ...string) error {
	command := createCommand(dir, printOutput, hide, cmdName, cmdArgs...)
	err := command.Start()
	if err != nil {
		log.Error("%s error", cmdName)
	} else {
		log.Debug("%s completed successfully", cmdName)
	}
	return err
}

func StopProcesses(image string) (string, error) {
	output, err := Exec("", false, nil, []int{128},
		"taskkill.exe", "/f", "/im", image,
	)
	return output, err
}

func FilePath(parts ...string) string {
	return strings.ReplaceAll(filepath.Join(parts...), "/", "\\")
}

func UncPath(server string, parts ...string) string {
	path := FilePath(server, filepath.Join(parts...))
	for !strings.HasPrefix(path, "\\\\") {
		path = "\\" + path
	}
	return path
}

func ProfilePath() string {
	return fsx.ForwardSlash(os.Getenv("USERPROFILE"))
}

func DownloadsPath() string {
	return filepath.Join(
		ProfilePath(),
		"Downloads",
	)
}

func Username() string {
	return os.Getenv("USERNAME")
}
