package fsx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golib/errs"
	"golib/lang"
	"golib/log"
)

const (
	DirPermissions  = 0755
	FilePermissions = 0644
)

func Split(path string) (stem string, extension string) {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[:i], path[i+1:]
		}
	}
	return path, ""
}

func AddSuffixBeforeExtension(path string, suffix string) string {
	stem, extension := Split(path)
	return stem + suffix + "." + extension
}

func EnsureUnique(template string) string {
	stem, extension := Split(template)
	return ensureUnique(stem, extension)
}

func EnsureUniqueWith(dir string, name string, extension string) string {
	return ensureUnique(filepath.Join(dir, name), extension)
}

func ensureUnique(stem string, extension string) string {
	candidate := fmt.Sprintf("%s.%s", stem, extension)
	if Exists(candidate) {
		var count = 'b'
		for {
			candidate = fmt.Sprintf("%s-%s.%s", stem, count, extension)
			if !Exists(candidate) {
				break
			}
			count++
			if count > 'z' {
				log.Debug("All filenames already taken for %s", candidate)
				os.Exit(2)
			}
		}
	}
	return candidate
}

func MustRead(name string) string {
	bytes, err := ReadString(name)
	lang.AssertNil(err)
	return string(bytes)
}

func WriteString(path string, content string) error {
	return WriteBytes(path, []byte(content))
}

func WriteBytes(path string, content []byte) error {
	return os.WriteFile(path, content, FilePermissions)
}

func ReadString(path string) (string, error) {
	bytes, err := ReadBytes(path)
	return string(bytes), err
}

func ReadBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ReadDir(path string) (results []os.FileInfo, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	var errors errs.Errors
	for _, entry := range entries {
		info, err := entry.Info()
		results = append(results, info)
		errors.Add(err)
	}
	err = errors.Combined()
	return
}

func List(path string, accept func(f os.DirEntry) bool) (files []string, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if accept(entry) {
			files = append(files, filepath.Join(path, entry.Name()))
		}
	}
	return
}

func Walk(path string, depth int, eachEntry func(dir string, f os.DirEntry)) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		eachEntry(path, entry)
	}
	if depth != 0 {
		for _, entry := range entries {
			if entry.IsDir() {
				err = Walk(filepath.Join(path, entry.Name()), depth-1, eachEntry)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

/*
	func MD5(path string) (hash string, err error) {
		f, err := os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()

		m := md5.New()
		_, err = io.Copy(m, f)
		if err != nil {
			return
		}

		sum := m.Sum(nil)
		hash = fmt.Sprintf("%x", sum)
		return
	}
*/

func Info(path string) os.FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}
	return info
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func SizeOf(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func WaitForFile(path string, notEmpty bool, delay time.Duration, timeout time.Duration) bool {
	time.Sleep(delay)
	return lang.WaitUntilTimeout(
		func() (time.Duration, bool) {
			if notEmpty {
				return delay, SizeOf(path) > 0
			} else {
				return delay, Exists(path)
			}
		},
		timeout,
	)
}

func Ensure(path string) error {
	log.Trace("MKDIR \"%s\"", path)
	err := os.MkdirAll(path, DirPermissions)
	if err != nil {
		log.Error("Unable to create [%s]: %s", path)
	}
	return err
}

func EnsureParent(path string) error {
	return Ensure(filepath.Dir(path))
}

func Rename(original string, renamed string) error {
	err := EnsureParent(renamed)
	if err != nil {
		log.Error("Could not move %s to %s: %s", original, renamed, err)
		return err
	}
	err = os.Rename(original, renamed)
	if err != nil {
		log.Error("Could not move %s to %s: %s", original, renamed, err)
		return err
	}
	log.Trace("MOVE  \"%s\" \"%s\"", original, renamed)
	return nil
}

/*
	func FileRename(path1 string, path2 string) error {
		if !Exists(path1) {
			return fmt.Errorf("Source file not found at [%s]", path1)
		}

		var backup string
		if Exists(path2) {
			backup := fmt.Sprintf("%s.%d.backup", path2, time.Now().UnixMilli())
			err := os.Rename(path2, backup)
			if err != nil {
				return err
			}
		}

		err := os.Rename(path1, path2)
		if err != nil {
			if len(backup) > 0 {
				// restore the original if possible
				os.Rename(backup, path2)
			}
			return err
		}

		return nil
	}
*/

func FindNamed(dir string, name string, depth int, stopOnError bool, onMatch func(path string) error) (errs lang.Errors) {
	return Find(dir,
		func(parent string, file os.DirEntry) (accept bool, stop bool) {
			return file.Name() == name, false
		},
		depth, stopOnError, onMatch)
}

func Find(dir string, accept func(parent string, file os.DirEntry) (accept bool, stop bool), depth int, stopOnError bool, onMatch func(path string) error) (errs lang.Errors) {
	log.Trace("CD    %s", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		errs.Set(dir, err)
		if stopOnError {
			return
		}
	}
	for _, file := range files {
		path := ForwardSlash(dir, file.Name())
		accepted, stop := accept(dir, file)
		if accepted {
			err = onMatch(path)
			if err != nil {
				errs.Set(path, err)
				if stopOnError {
					return
				}
			}
		}
		if stop {
			return
		}
		if depth > 1 && file.IsDir() {
			errs2 := Find(path, accept, depth-1, stopOnError, onMatch)
			if errs2.Found() {
				errs.AddAll(errs2)
				if stopOnError {
					return
				}
			}
		}
	}
	return
}

func ForwardSlash(parts ...string) string {
	return strings.ReplaceAll(filepath.Join(parts...), "\\", "/")
}
