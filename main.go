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

	SlideFormat = make(map[int][]Drawn)

	outPath := filepath.Join(rootPath, "slides")
	os.MkdirAll(outPath, 0777)

	runtime.LockOSThread()

	go func() {
		for {
			inStr := <-InTPickerChannel
			TextFromTPicker = pickText(inStr)
			ClearAfterTPicker = true
		}
	}()

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

		if ClearAfterTPicker {
			if len(TextFromTPicker) > 0 {
				// if isUpdateDialog {
				// 	oldComment := comments[commentIndexBeingUpdated]
				// 	comments[commentIndexBeingUpdated] = Comment{oldComment.X, oldComment.Y, TextFromTPicker}
				// 	isUpdateDialog = false
				// 	commentIndexBeingUpdated = 0
				// } else {
				// 	comments = append(comments, Comment{activeX, activeY, TextFromTPicker})
				// }
				textDetail := TextDetail{TextFromTPicker, currentTextColor, 1}
				TextDetails = append(TextDetails, textDetail)
				drawn := Drawn{Type: TextType, X: activeX, Y: activeY, DetailsId: len(TextDetails) - 1}
				objs := SlideFormat[CurrentSlide]
				SlideFormat[CurrentSlide] = append(objs, drawn)
			}
			TextFromTPicker = ""
			activeX, activeY = -1, -1

			DrawWorkView(window, CurrentSlide)
			window.SetMouseButtonCallback(workViewMouseCallback)
			// window.SetCursorPosCallback(getHoverCB(objCoords))

			ClearAfterTPicker = false
		}

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}

}
