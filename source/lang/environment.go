package lang

import (
	"os"
	"time"
)

func ProgramExeDate() *time.Time {
	exe, err := os.Executable()
	if err == nil {
		file, err := os.Stat(exe)
		if err == nil {
			f := file.ModTime()
			return &f
		}
	}
	return nil
}
