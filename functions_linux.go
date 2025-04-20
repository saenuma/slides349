package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetExecPath(execName string) string {
	homeDir, _ := os.UserHomeDir()
	var cmdPath string
	begin := os.Getenv("SNAP")
	cmdPath = filepath.Join(homeDir, "bin", execName)
	if begin != "" && !strings.HasPrefix(begin, "/snap/go/") {
		cmdPath = filepath.Join(begin, "bin", execName)
	}

	return cmdPath
}

func pickFileUbuntu(exts string) string {
	fPickerPath := GetExecPath("fpicker")

	rootPath, _ := GetRootPath()
	cmd := exec.Command(fPickerPath, rootPath, exts)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}

func PickImageFile() string {
	return pickFileUbuntu("png|jpg")
}

func pickText(inStr string) string {
	execPath := GetExecPath("tpicker")
	cmd := exec.Command(execPath)
	if inStr != "" {
		cmd = exec.Command(execPath, inStr)
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}

func pickColor() string {
	aCPickerPath := GetExecPath("acpicker")

	cmd := exec.Command(aCPickerPath)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}
