package main

import (
	"fmt"
	"image"
	"math"
	"slices"
	"strconv"
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/kovidgoyal/imaging"
)

func DrawWorkView(window *glfw.Window, slide int) {
	CurrentSlide = slide

	window.SetTitle(fmt.Sprintf("Project: %s ---- %s", ProjectName, ProgTitle))

	ObjCoords = make(map[int]g143.Rect)

	wWidth, wHeight := window.GetSize()
	theCtx := New2dCtx(wWidth, wHeight, &ObjCoords)

	// top panel
	sTRS := theCtx.drawButtonA(AddSlideBtn, 50, 10, "Add Slide", fontColor, "#D9D5B0", "#D9D5B0")
	tTX := nextHorizontalX(sTRS, 20)
	tTRS := theCtx.drawButtonA(TextTool, tTX, 10, "Text", fontColor, "#D9D5B0", "#D9D5B0")
	iTX := nextHorizontalX(tTRS, 20)
	iTRS := theCtx.drawButtonA(ImageTool, iTX, 10, "Image", fontColor, "#D9D5B0", "#D9D5B0")
	mTX := nextHorizontalX(iTRS, 20)
	mTRS := theCtx.drawButtonA(MoveTool, mTX, 10, "Move", fontColor, "#D9D5B0", "#D9D5B0")
	dividerX := nextHorizontalX(mTRS, 20)
	theCtx.ggCtx.SetHexColor("#444")
	theCtx.ggCtx.DrawRectangle(float64(dividerX), 10, 2, float64(tTRS.Height))
	theCtx.ggCtx.Fill()

	// pencil extras
	mSRS := theCtx.drawButtonB(MinusSizeBtn, dividerX+20, 10+5, "--", "#fff", "#aaa")
	tSIX := nextHorizontalX(mSRS, 5)

	size := InputsState["size"]
	tSIRS := theCtx.drawInput(DrawnSizeInput, tSIX, 10, size)
	pSX := nextHorizontalX(tSIRS, 5)
	pSRS := theCtx.drawButtonB(PlusSizeBtn, pSX, 10+5, "+", "#fff", "#aaa")
	tCX := nextHorizontalX(pSRS, 20)
	selectedColor := InputsState["color"]
	theCtx.drawColorBox(ColorPickerBtn, tCX, 10+5, tTRS.Height-15, selectedColor)

	// place indicator on activeTool
	activeToolRS := ObjCoords[activeTool]
	theCtx.drawButtonA(activeTool, activeToolRS.OriginX, activeToolRS.OriginY, toolNames[activeTool],
		fontColor, "#D9D5B0", "#AEAC9C")

	// slides panel
	slideWidth, slideHeight := int(math.Ceil(WorkAreaWidth*0.8))-2, int(math.Ceil(WorkAreaHeight*0.8))-2
	currentY := 80
	for i := range TotalSlides {
		displayI := i + 1
		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawString(fmt.Sprintf("%d", displayI), 10, float64(currentY+FontSize))

		// draw border rectangle
		theCtx.ggCtx.SetHexColor(fontColor)
		slideX := 10 + FontSize + 5
		theCtx.ggCtx.DrawRectangle(float64(slideX), float64(currentY), WorkAreaWidth*0.15, WorkAreaHeight*0.15+1)
		theCtx.ggCtx.Fill()

		// draw indicator for current slide
		if i == CurrentSlide {
			theCtx.ggCtx.SetHexColor("#BE7171")
			theCtx.ggCtx.DrawRectangle(float64(slideX)+WorkAreaWidth*0.15+4, float64(currentY), 10, WorkAreaHeight*0.15)
			theCtx.ggCtx.Fill()
		}

		// draw thumbnail
		slideImg := drawSlide(i, slideWidth, slideHeight)
		thumbnailWidth, thumbnailHeight := int(math.Ceil(WorkAreaWidth*0.15))-2, int(math.Ceil(WorkAreaHeight*0.15))-2
		slideImg = imaging.Fit(slideImg, thumbnailWidth, thumbnailHeight, imaging.Lanczos)
		theCtx.ggCtx.DrawImage(slideImg, slideX+1, currentY+1)

		// register thumbnail
		sTW, sTH := int(math.Ceil(WorkAreaWidth*0.15)), int(math.Ceil(WorkAreaHeight*0.15))
		SlideThumbnailRect := g143.NewRect(slideX, currentY, sTW, sTH)
		ObjCoords[1000+displayI] = SlideThumbnailRect

		currentY += 10 + int(math.Ceil(WorkAreaHeight*0.15))
	}

	// canvas
	canvasX := int(math.Ceil(WorkAreaWidth*0.15)) + 10 + FontSize + 5 + 30

	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawRectangle(float64(canvasX), 80, WorkAreaWidth*0.8, WorkAreaHeight*0.8)
	theCtx.ggCtx.Fill()

	CanvasRect = g143.NewRect(canvasX+1, 80+1, slideWidth, slideHeight)
	ObjCoords[CanvasWidget] = CanvasRect

	canvasImg := drawSlide(CurrentSlide, slideWidth, slideHeight)
	theCtx.ggCtx.DrawImage(canvasImg, canvasX+1, 80+1)

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}

func drawSlide(slideNo int, workingWidth, workingHeight int) image.Image {
	// frame buffer
	ggCtx := gg.NewContext(workingWidth, workingHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(workingWidth), float64(workingHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	canvasRS := ObjCoords[CanvasWidget]
	// load font
	fontPath := GetDefaultFontPath()
	ggCtx.LoadFontFace(fontPath, 20)

	for i, obj := range SlideFormat[slideNo] {
		if obj.Type == TextType {
			strs := strings.Split(strings.ReplaceAll(obj.Text, "\r", ""), "\n")
			textFontSize := float64(obj.Size) * 15
			textFontSizeInt := int(math.Ceil(textFontSize))
			ggCtx.LoadFontFace(fontPath, textFontSize)

			maxX := 0
			currentY := obj.Y

			for _, str := range strs {
				strW, _ := ggCtx.MeasureString(str)
				if int(strW) > maxX {
					maxX = int(strW)
				}
				ggCtx.SetHexColor(obj.Color)
				drawnX := obj.X - canvasRS.OriginX
				drawnY := currentY + textFontSizeInt - canvasRS.OriginY
				ggCtx.DrawString(str, float64(drawnX), float64(drawnY))
				currentY += 10 + textFontSizeInt
			}

			obj.W = maxX
			obj.H = currentY - obj.Y
			SlideFormat[slideNo][i] = obj

		} else if obj.Type == ImageType {
			img, err := imaging.Open(obj.ImagePath, imaging.AutoOrientation(true))
			if err != nil {
				fmt.Println(err)
			}

			imgW := float64(obj.Size) * 100
			scale := float64(imgW) / float64(img.Bounds().Dx())
			newH := int(scale * float64(img.Bounds().Dy()))
			img = imaging.Fit(img, int(imgW), newH, imaging.Lanczos)
			ggCtx.DrawImage(img, obj.X-canvasRS.OriginX, obj.Y-canvasRS.OriginY)

			obj.W = int(imgW)
			obj.H = newH

			SlideFormat[slideNo][i] = obj
		}
	}

	return ggCtx.Image()
}

func workViewMouseCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	wWidth, wHeight := window.GetSize()

	var widgetCode int

	for code, RS := range ObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	// rootPath, _ := GetRootPath()

	switch widgetCode {
	case AddSlideBtn:
		TotalSlides += 1
		CurrentSlide += 1
		SlideFormat = slices.Insert(SlideFormat, CurrentSlide, make([]Drawn, 0))
		DrawWorkView(window, CurrentSlide)

	case TextTool, ImageTool, MoveTool:
		activeTool = widgetCode

		theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)

		// clear all tools
		for _, toolId := range []int{TextTool, ImageTool, MoveTool} {
			toolRS := ObjCoords[toolId]
			theCtx.drawButtonA(toolId, toolRS.OriginX, toolRS.OriginY, toolNames[toolId],
				fontColor, "#D9D5B0", "#D9D5B0")
		}
		// place indicator on activeTool
		activeToolRS := ObjCoords[activeTool]
		theCtx.drawButtonA(activeTool, activeToolRS.OriginX, activeToolRS.OriginY, toolNames[activeTool],
			fontColor, "#D9D5B0", "#AEAC9C")

		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		CurrentWindowFrame = theCtx.ggCtx.Image()

	case CanvasWidget:
		// theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)
		ctrlState := window.GetKey(glfw.KeyLeftControl)
		// canvasRS := ObjCoords[CanvasWidget]

		// translastedMouseX, translatedMouseY := xPos-float64(canvasRS.OriginX), yPos-float64(canvasRS.OriginY)
		activeX, activeY = xPosInt, yPosInt

		if activeTool == TextTool {

			if ctrlState == glfw.Release {
				// stop interaction till returning from tpicker
				window.SetMouseButtonCallback(nil)
				window.SetCursorPosCallback(nil)

				foundIndex := -1
				for i, obj := range SlideFormat[CurrentSlide] {
					if obj.Type != TextType {
						continue
					}
					objRect := g143.NewRect(obj.X, obj.Y, obj.W, obj.H)
					if g143.InRect(objRect, activeX, activeY) {
						foundIndex = i
						break
					}
				}

				if foundIndex != -1 {
					DrawnEditIndex = foundIndex
					PickerChan <- []string{"text", SlideFormat[CurrentSlide][foundIndex].Text}
				} else {
					PickerChan <- []string{"text", ""}
				}

			} else if ctrlState == glfw.Press {
				foundIndex := -1
				for i, obj := range SlideFormat[CurrentSlide] {
					if obj.Type != TextType {
						continue
					}
					objRect := g143.NewRect(obj.X, obj.Y, obj.W, obj.H)
					if g143.InRect(objRect, activeX, activeY) {
						foundIndex = i
						break
					}
				}

				if foundIndex != -1 {
					objs := SlideFormat[CurrentSlide]
					SlideFormat[CurrentSlide] = slices.Delete(objs, foundIndex, foundIndex+1)
				}

				DrawWorkView(window, CurrentSlide)
			}

		} else if activeTool == ImageTool {

			if ctrlState == glfw.Release {
				// stop interaction till returning from tpicker
				window.SetMouseButtonCallback(nil)
				window.SetCursorPosCallback(nil)

				foundIndex := -1
				for i, obj := range SlideFormat[CurrentSlide] {
					if obj.Type != ImageType {
						continue
					}
					objRect := g143.NewRect(obj.X, obj.Y, obj.W, obj.H)
					if g143.InRect(objRect, activeX, activeY) {
						foundIndex = i
						break
					}
				}

				if foundIndex != -1 {
					DrawnEditIndex = foundIndex
					PickerChan <- []string{"image", SlideFormat[CurrentSlide][foundIndex].ImagePath}
				} else {
					PickerChan <- []string{"image", ""}
				}

			} else if ctrlState == glfw.Press {
				foundIndex := -1
				for i, obj := range SlideFormat[CurrentSlide] {
					if obj.Type != ImageType {
						continue
					}
					objRect := g143.NewRect(obj.X, obj.Y, obj.W, obj.H)
					if g143.InRect(objRect, activeX, activeY) {
						foundIndex = i
						break
					}
				}

				if foundIndex != -1 {
					objs := SlideFormat[CurrentSlide]
					SlideFormat[CurrentSlide] = slices.Delete(objs, foundIndex, foundIndex+1)
				}

				DrawWorkView(window, CurrentSlide)
			}

		} else if activeTool == MoveTool {
			if ToMoveIndex == -1 {
				foundIndex := -1
				for i, obj := range SlideFormat[CurrentSlide] {
					objRect := g143.NewRect(obj.X, obj.Y, obj.W, obj.H)
					if g143.InRect(objRect, activeX, activeY) {
						foundIndex = i
						break
					}
				}

				ToMoveIndex = foundIndex

			} else {
				obj := SlideFormat[CurrentSlide][ToMoveIndex]
				obj.X = activeX
				obj.Y = activeY

				SlideFormat[CurrentSlide][ToMoveIndex] = obj
				ToMoveIndex = -1

				DrawWorkView(window, CurrentSlide)
			}

		}

	case PlusSizeBtn:
		size := InputsState["size"]
		sizeInt, _ := strconv.Atoi(size)
		if sizeInt != 10 {
			sizeInt += 1
		}

		theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)
		widgetRS := ObjCoords[DrawnSizeInput]
		theCtx.drawInput(DrawnSizeInput, widgetRS.OriginX, widgetRS.OriginY, strconv.Itoa(sizeInt))

		InputsState["size"] = strconv.Itoa(sizeInt)
		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		CurrentWindowFrame = theCtx.ggCtx.Image()

	case MinusSizeBtn:
		size := InputsState["size"]
		sizeInt, _ := strconv.Atoi(size)
		if sizeInt != 1 {
			sizeInt -= 1
		}
		theCtx := Continue2dCtx(CurrentWindowFrame, &ObjCoords)
		widgetRS := ObjCoords[DrawnSizeInput]
		theCtx.drawInput(DrawnSizeInput, widgetRS.OriginX, widgetRS.OriginY, strconv.Itoa(sizeInt))

		InputsState["size"] = strconv.Itoa(sizeInt)
		// send the frame to glfw window
		g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
		window.SwapBuffers()

		// save the frame
		CurrentWindowFrame = theCtx.ggCtx.Image()

	case ColorPickerBtn:
		window.SetMouseButtonCallback(nil)
		window.SetCursorPosCallback(nil)
		PickerChan <- []string{"color", ""}
	}

	// for generated buttons
	if widgetCode > 1000 && widgetCode < 2000 {
		slideNum := widgetCode - 1000 - 1
		ctrlState := window.GetKey(glfw.KeyLeftControl)

		if ctrlState == glfw.Release {
			CurrentSlide = slideNum
			DrawWorkView(window, CurrentSlide)

			ToMoveIndex = -1
		} else if ctrlState == glfw.Press {
			if TotalSlides != 1 {
				SlideFormat = slices.Delete(SlideFormat, slideNum-1, slideNum)
				CurrentSlide = 0
				TotalSlides -= 1
				DrawWorkView(window, CurrentSlide)
				ToMoveIndex = -1
			}
		}

	}

}
