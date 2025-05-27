package main

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/pkg/errors"
)

func GetDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "v349_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}

func GetRootPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "os error")
	}

	dd := os.Getenv("SNAP_USER_COMMON")

	if strings.HasPrefix(dd, filepath.Join(hd, "snap", "go")) || dd == "" {
		dd = filepath.Join(hd, "Videos349")
	}

	retPath := filepath.Join(dd, "slides")
	os.MkdirAll(retPath, 0777)

	return retPath, nil
}

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func UntestedRandomString(length int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetProjectFiles() []ToSortProject {
	// display some project names
	rootPath, _ := GetRootPath()
	dirEs, _ := os.ReadDir(rootPath)

	projectFiles := make([]ToSortProject, 0)
	for _, dirE := range dirEs {
		if dirE.IsDir() {
			continue
		}

		if strings.HasSuffix(dirE.Name(), ".v3p") {
			fInfo, _ := dirE.Info()
			projectFiles = append(projectFiles, ToSortProject{dirE.Name(), fInfo.ModTime()})
		}
	}

	slices.SortFunc(projectFiles, func(a, b ToSortProject) int {
		return b.ModTime.Compare(a.ModTime)
	})

	return projectFiles
}

func ExternalLaunch(p string) {
	cmd := "url.dll,FileProtocolHandler"
	runDll32 := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")

	if runtime.GOOS == "windows" {
		exec.Command(runDll32, cmd, p).Run()
	} else if runtime.GOOS == "linux" {
		exec.Command("xdg-open", p).Run()
	}
}
