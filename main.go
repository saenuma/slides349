package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	rootPath, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	// initialize
	SlideFormat = make(map[int][]Drawn)
	SlideMemory = make(map[int]map[string]string)
	SlideMemory[1] = map[string]string{
		"color": "#444",
		"size":  "1",
	}

	outPath := filepath.Join(rootPath, "slides")
	os.MkdirAll(outPath, 0777)

	runtime.LockOSThread()

	go func() {
		for {
			aSlice := <-PickerChan
			if aSlice[0] == "text" {
				TextFromTPicker = pickText(aSlice[1])
				ClearAfterTPicker = true
			} else if aSlice[0] == "color" {
				TextFromACPicker = pickColor()
				ClearAFterACPicker = true
			}
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
	// quick hover effect
	window.SetCursorPosCallback(getHoverCB(&ObjCoords))

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		if ClearAfterTPicker {
			if len(TextFromTPicker) > 0 {
				size := SlideMemory[CurrentSlide]["size"]
				sizeInt, _ := strconv.Atoi(size)

				if DrawnEditIndex != -1 {
					drawnText := SlideFormat[CurrentSlide][DrawnEditIndex]
					drawnText.Text = TextFromTPicker
					SlideFormat[CurrentSlide][DrawnEditIndex] = drawnText
					DrawnEditIndex = -1
				} else {
					drawn := Drawn{Type: TextType, X: activeX, Y: activeY, Text: TextFromTPicker,
						Color: SlideMemory[CurrentSlide]["color"], Size: sizeInt}
					objs := SlideFormat[CurrentSlide]
					SlideFormat[CurrentSlide] = append(objs, drawn)
				}
			}
			TextFromTPicker = ""
			activeX, activeY = -1, -1

			DrawWorkView(window, CurrentSlide)
			window.SetMouseButtonCallback(workViewMouseCallback)
			window.SetCursorPosCallback(getHoverCB(&ObjCoords))

			ClearAfterTPicker = false
		}

		if ClearAFterACPicker {
			SlideMemory[CurrentSlide]["color"] = TextFromACPicker
			TextFromACPicker = ""

			DrawWorkView(window, CurrentSlide)
			window.SetMouseButtonCallback(workViewMouseCallback)
			window.SetCursorPosCallback(getHoverCB(&ObjCoords))

			ClearAFterACPicker = false
		}

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}

}
