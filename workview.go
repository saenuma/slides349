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
	InWorkView = true
	window.SetTitle(fmt.Sprintf("Project: %s ---- %s", ProjectName, ProgTitle))

	ObjCoords = make(map[int]g143.Rect)

	wWidth, wHeight := window.GetSize()
	theCtx := New2dCtx(wWidth, wHeight, &ObjCoords)

	// top panel
	bBRect := theCtx.drawButtonA(BackBtn, 50, 10, "Back", "#fff", "#845B5B", "#845B5B")
	aTX := nextHorizontalX(bBRect, 40)
	aSRS := theCtx.drawButtonA(AddSlideBtn, aTX, 10, "Add Slide", fontColor, "#D9D5B0", "#D9D5B0")
	tTX := nextHorizontalX(aSRS, 20)
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
	pSX := nextHorizontalX(tSIRS, 15)
	pSRS := theCtx.drawButtonB(PlusSizeBtn, pSX, 10+5, "+", "#fff", "#aaa")
	tCX := nextHorizontalX(pSRS, 20)
	selectedColor := InputsState["color"]
	cPBRS := theCtx.drawColorBox(ColorPickerBtn, tCX, 10+5, tTRS.Height-15, selectedColor)
	fCX := nextHorizontalX(cPBRS, 20)
	selectedFontClass := InputsState["font"]
	theCtx.drawInput(FontNameInput, fCX, 10, selectedFontClass)

	// place indicator on activeTool
	activeToolRS := ObjCoords[activeTool]
	theCtx.drawButtonA(activeTool, activeToolRS.OriginX, activeToolRS.OriginY, toolNames[activeTool],
		fontColor, "#D9D5B0", "#AEAC9C")

	genRange := func(a, b, total int) []int {
		ret := make([]int, 0)
		count := 0
		for i := a; i < b; i++ {
			if count == 5 {
				break
			}
			if i == total {
				break
			}
			ret = append(ret, i)
			count += 1
		}
		return ret
	}

	// get number of slides to display in GUI
	beginIndex := SlidesOffset * 5
	endIndex := (SlidesOffset + 1) * 5
	slides := genRange(beginIndex, endIndex, TotalSlides)
	// slides panel
	slideWidth, slideHeight := int(math.Ceil(WorkAreaWidth*0.8))-2, int(math.Ceil(WorkAreaHeight*0.8))-2
	currentY := 80
	slideX := 10 + FontSize + 5
	// canvas
	canvasX := int(math.Ceil(WorkAreaWidth*0.15)) + 10 + FontSize + 5 + 30

	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawRectangle(float64(canvasX), 80, WorkAreaWidth*0.8, WorkAreaHeight*0.8)
	theCtx.ggCtx.Fill()

	CanvasRect = g143.NewRect(canvasX+1, 80+1, slideWidth, slideHeight)
	ObjCoords[CanvasWidget] = CanvasRect

	canvasImg := drawSlide(CurrentSlide, true)
	theCtx.ggCtx.DrawImage(canvasImg, canvasX+1, 80+1)

	for _, i := range slides {
		displayI := i + 1
		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawString(fmt.Sprintf("%d", displayI), 10, float64(currentY+FontSize))

		// draw border rectangle
		theCtx.ggCtx.SetHexColor(fontColor)
		theCtx.ggCtx.DrawRectangle(float64(slideX), float64(currentY), WorkAreaWidth*0.15, WorkAreaHeight*0.15+1)
		theCtx.ggCtx.Fill()

		// draw indicator for current slide
		if i == CurrentSlide {
			theCtx.ggCtx.SetHexColor("#BE7171")
			theCtx.ggCtx.DrawRectangle(float64(slideX)+WorkAreaWidth*0.15+4, float64(currentY), 10, WorkAreaHeight*0.15)
			theCtx.ggCtx.Fill()
		}

		// draw thumbnail
		slideImg := drawSlide(i, false)
		thumbnailWidth, thumbnailHeight := int(math.Ceil(WorkAreaWidth*0.15))-2, int(math.Ceil(WorkAreaHeight*0.15))-2
		thumbImg := imaging.Thumbnail(slideImg, thumbnailWidth, thumbnailHeight, imaging.Lanczos)
		theCtx.ggCtx.DrawImage(thumbImg, slideX+1, currentY+1)

		// register thumbnail
		sTW, sTH := int(math.Ceil(WorkAreaWidth*0.15)), int(math.Ceil(WorkAreaHeight*0.15))
		SlideThumbnailRect := g143.NewRect(slideX, currentY, sTW, sTH)
		ObjCoords[1000+displayI] = SlideThumbnailRect

		currentY += 10 + int(math.Ceil(WorkAreaHeight*0.15))
	}

	// write totalSlides
	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawString(fmt.Sprintf("Total Slides: %d", TotalSlides), float64(slideX),
		float64(slideHeight)+FontSize+30+float64(CanvasRect.OriginY))

	// send the frame to glfw window
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), theCtx.windowRect())
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}

func drawSlide(slideNo int, forGui bool) image.Image {

	var workingWidth, workingHeight int
	var scale float64
	if forGui {
		workingWidth, workingHeight = int(math.Ceil(WorkAreaWidth*0.8))-2, int(math.Ceil(WorkAreaHeight*0.8))-2
		scale = 1.0
	} else {
		workingWidth, workingHeight = WorkAreaWidth, WorkAreaHeight
		scale = float64(workingWidth) / ((WorkAreaWidth * 0.8) - 2)
	}

	// frame buffer
	ggCtx := gg.NewContext(workingWidth, workingHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(workingWidth), float64(workingHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	canvasRS := ObjCoords[CanvasWidget]
	for i, obj := range SlideFormat[slideNo] {
		if obj.Type == TextType {
			strs := strings.Split(strings.ReplaceAll(obj.Text, "\r", ""), "\n")
			textFontSize := float64(obj.Size) * 15 * scale
			fontName := "regular"
			if obj.FontName != "" {
				fontName = obj.FontName
			}
			fontDef := getFontDef(fontName)
			fontPath := getFontPath(fontDef.Index)
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

			imgW := float64(obj.Size) * 100 * scale
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
	case BackBtn:
		SaveSlideProject()

		// clear some variables
		SlideFormat = make([][]Drawn, 0)
		ProjectName = ""
		window.SetTitle(ProgTitle)

		// redraw
		DrawBeginView(window)
		window.SetMouseButtonCallback(projViewMouseCallback)
		window.SetKeyCallback(ProjKeyCallback)
		window.SetCursorPosCallback(getHoverCB(&ProjObjCoords))
		window.SetScrollCallback(nil)

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
		if sizeInt != 5 {
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

	case FontNameInput:
		drawFontPicker(window, CurrentWindowFrame)
		window.SetMouseButtonCallback(fontPickerMouseCallback)
		window.SetCursorPosCallback(getHoverCB(&FDObjCoords))
		window.SetScrollCallback(nil)
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
				SlideFormat = slices.Delete(SlideFormat, slideNum, slideNum+1)
				CurrentSlide = 0
				TotalSlides -= 1
				DrawWorkView(window, CurrentSlide)
				ToMoveIndex = -1
			}
		}

	}

}

func workViewScrollCB(window *glfw.Window, xoff float64, yoff float64) {
	// makes the scroll more realistic
	scrollEventsCount += 1
	if scrollEventsCount != 5 {
		return
	} else {
		scrollEventsCount = 0
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// 	wWidth, wHeight := window.GetSize()
	slidesPanelRect := g143.NewRect(10, CanvasRect.OriginY, int(math.Ceil(WorkAreaWidth*0.15)),
		CanvasRect.Height)

	if g143.InRect(slidesPanelRect, xPosInt, yPosInt) {
		if int(yoff) == -1 {
			// show bottom slides
			if ((SlidesOffset + 1) * 5) < TotalSlides {
				SlidesOffset += 1
				CurrentSlide = SlidesOffset * 5
				DrawWorkView(window, CurrentSlide)
			}
		} else if int(yoff) == 1 {
			// show top slides
			if SlidesOffset-1 >= 0 {
				SlidesOffset -= 1
				CurrentSlide = SlidesOffset * 5
				DrawWorkView(window, CurrentSlide)
			}
		}
	}

}
