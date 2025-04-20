package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sqweek/dialog"
)

func PickImageFile() string {
	filename, err := dialog.File().Filter("PNG Image", "png").Filter("JPEG Image", "jpg").Load()
	if filename == "" || err != nil {
		log.Println(err)
		return ""
	}
	return filename
}

func pickText(inStr string) string {
	execPath, _ := os.Executable()
	cmdPath := filepath.Join(filepath.Dir(execPath), "tpicker.exe")

	cmd := exec.Command(cmdPath)
	if inStr != "" {
		cmd = exec.Command(cmdPath, inStr)
	}

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}

func pickColor() string {
	execPath, _ := os.Executable()
	cmdPath := filepath.Join(filepath.Dir(execPath), "acpicker.exe")
	cmd := exec.Command(cmdPath)

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return strings.TrimSpace(string(out))
}
