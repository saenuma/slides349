package main

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	outPath := filepath.Join(rootPath, "slides")
	os.MkdirAll(outPath, 0777)

	runtime.LockOSThread()

	window := g143.NewWindow(1400, 800, ProgTitle, false)
	DrawWorkView(window, 1)

	// respond to the mouse
	window.SetMouseButtonCallback(workViewMouseCallback)
	// // respond to the keyboard
	// window.SetKeyCallback(ProjKeyCallback)
	// // save the project file
	// window.SetCloseCallback(SaveProjectCloseCallback)
	// // quick hover effect
	// window.SetCursorPosCallback(getHoverCB(ProjObjCoords))

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}

}
