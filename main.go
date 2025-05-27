package main

import (
	"runtime"
	"strconv"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	_, err := GetRootPath()
	if err != nil {
		panic(err)
	}

	InputsState = make(map[string]string)
	InputsState = map[string]string{
		"color": "#444",
		"size":  "1",
	}
	activeTool = TextTool

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
			} else if aSlice[0] == "image" {
				tmp := PickImageFile()
				PathFromFPicker = ""
				if tmp == "" && len(aSlice[1]) != 0 {
					PathFromFPicker = aSlice[1]
				} else {
					PathFromFPicker = tmp
				}
				ClearAfterFPicker = true
			}
		}
	}()

	window := g143.NewWindow(1400, 800, ProgTitle, false)

	DrawBeginView(window)
	// respond to the mouse
	window.SetMouseButtonCallback(projViewMouseCallback)
	// respond to the keyboard
	window.SetKeyCallback(ProjKeyCallback)
	// save the project file
	window.SetCloseCallback(func(w *glfw.Window) {
		SaveSlideProject()
	})
	// quick hover effect
	window.SetCursorPosCallback(getHoverCB(&ProjObjCoords))

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		if ClearAfterTPicker {
			if len(TextFromTPicker) > 0 {
				size := InputsState["size"]
				sizeInt, _ := strconv.Atoi(size)

				if DrawnEditIndex != -1 {
					drawnText := SlideFormat[CurrentSlide][DrawnEditIndex]
					drawnText.Text = TextFromTPicker
					drawnText.Size = sizeInt
					drawnText.Color = InputsState["color"]
					SlideFormat[CurrentSlide][DrawnEditIndex] = drawnText
					DrawnEditIndex = -1
				} else {
					objs := SlideFormat[CurrentSlide]
					toWriteWidgetCode := 8001
					if len(objs) > 0 {
						toWriteWidgetCode = objs[len(objs)-1].WidgetCode + 1
					}
					drawn := Drawn{Type: TextType, X: activeX, Y: activeY, Text: TextFromTPicker,
						Color: InputsState["color"], Size: sizeInt, WidgetCode: toWriteWidgetCode}

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
			InputsState["color"] = TextFromACPicker
			TextFromACPicker = ""

			DrawWorkView(window, CurrentSlide)
			window.SetMouseButtonCallback(workViewMouseCallback)
			window.SetCursorPosCallback(getHoverCB(&ObjCoords))

			ClearAFterACPicker = false
		}

		if ClearAfterFPicker {
			if PathFromFPicker != "" {
				objs := SlideFormat[CurrentSlide]
				size := InputsState["size"]
				sizeInt, _ := strconv.Atoi(size)

				if DrawnEditIndex != -1 {
					drawnImg := SlideFormat[CurrentSlide][DrawnEditIndex]
					drawnImg.ImagePath = PathFromFPicker
					drawnImg.Size = sizeInt
					SlideFormat[CurrentSlide][DrawnEditIndex] = drawnImg
					DrawnEditIndex = -1
				} else {
					toWriteWidgetCode := 8001
					if len(objs) > 0 {
						toWriteWidgetCode = objs[len(objs)-1].WidgetCode + 1
					}
					drawn := Drawn{Type: ImageType, X: activeX, Y: activeY, ImagePath: PathFromFPicker,
						Size: sizeInt, WidgetCode: toWriteWidgetCode}
					SlideFormat[CurrentSlide] = append(objs, drawn)
				}
			}
			PathFromFPicker = ""
			activeX, activeY = -1, -1

			DrawWorkView(window, CurrentSlide)
			window.SetMouseButtonCallback(workViewMouseCallback)
			window.SetCursorPosCallback(getHoverCB(&ObjCoords))

			ClearAfterFPicker = false
		}

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}

}
