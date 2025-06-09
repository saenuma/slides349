package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func DrawBeginView(window *glfw.Window) {
	ProjObjCoords = make(map[int]g143.Rect)
	wWidth, wHeight := window.GetSize()

	theCtx := New2dCtx(wWidth, wHeight, &ProjObjCoords)

	fontPath := getFontPath(1)
	theCtx.ggCtx.LoadFontFace(fontPath, 30)

	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawString("New Project", 20, 10+30)

	theCtx.ggCtx.LoadFontFace(fontPath, 20)
	pnIRect := theCtx.drawInputB(PROJ_NameInput, 20, 60, 420, "enter project name", false)
	pnBtnX := nextHorizontalX(pnIRect, 30)
	nPRS := theCtx.drawButtonB(PROJ_NewProject, pnBtnX, 60, "New Project", fontColor, "#B3AE97")
	oWDBX := nextHorizontalX(nPRS, 40)
	theCtx.drawButtonA(PROJ_OpenWDBtn, oWDBX, 10, "Open Working Directory", "#fff", "#56845A", "#56845A")

	// second row border
	borderY := nextVerticalY(pnIRect, 10)
	theCtx.ggCtx.SetHexColor("#999")
	theCtx.ggCtx.DrawRoundedRectangle(10, float64(borderY), float64(wWidth)-20, 2, 2)
	theCtx.ggCtx.Fill()

	theCtx.ggCtx.LoadFontFace(fontPath, 30)
	theCtx.ggCtx.SetHexColor(fontColor)
	theCtx.ggCtx.DrawString("Continue Projects", 20, float64(borderY)+12+30)
	theCtx.ggCtx.LoadFontFace(fontPath, 20)

	projectFiles := GetProjectFiles()
	currentX := 40
	currentY := borderY + 22 + 30 + 10
	for i, pf := range projectFiles {

		btnId := 1000 + (i + 1)
		pfRect := theCtx.drawButtonB(btnId, currentX, currentY, pf.Name, "#fff", "#5F699F")

		newX := currentX + pfRect.Width + 10
		if newX > (wWidth - pfRect.Width) {
			currentY += 50
			currentX = 40
		} else {
			currentX += pfRect.Width + 10
		}

	}

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}

func projViewMouseCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.Rect
	var widgetCode int

	for code, RS := range ProjObjCoords {
		if g143.InRect(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	rootPath, _ := GetRootPath()

	switch widgetCode {
	case PROJ_NewProject:
		if NameInputEnteredTxt == "" {
			return
		}

		// initialize
		SlideFormat = make([][]Drawn, 0)
		SlideFormat = append(SlideFormat, make([]Drawn, 0)) // first slide

		// create file
		ProjectName = NameInputEnteredTxt + ".s3p"
		outPath := filepath.Join(rootPath, ProjectName)
		os.WriteFile(outPath, []byte(""), 0777)

		// move to work view
		DrawWorkView(window, CurrentSlide)
		window.SetMouseButtonCallback(workViewMouseCallback)
		window.SetKeyCallback(nil)
		window.SetScrollCallback(workViewScrollCB)
		// quick hover effect
		window.SetCursorPosCallback(getHoverCB(&ObjCoords))

	case PROJ_OpenWDBtn:
		rootPath, _ := GetRootPath()
		ExternalLaunch(rootPath)
	}

	if widgetCode > 1000 && widgetCode < 2000 {
		num := widgetCode - 1000 - 1
		projectFile := GetProjectFiles()[num]

		ProjectName = projectFile.Name

		// load slide data
		var objs [][]Drawn
		rootPath, _ := GetRootPath()
		inPath := filepath.Join(rootPath, ProjectName)
		rawBytes, _ := os.ReadFile(inPath)
		json.Unmarshal(rawBytes, &objs)

		TotalSlides = len(objs)
		SlideFormat = objs

		// move to work view
		DrawWorkView(window, CurrentSlide)
		window.SetMouseButtonCallback(workViewMouseCallback)
		window.SetKeyCallback(nil)
		window.SetScrollCallback(workViewScrollCB)
		window.SetCursorPosCallback(getHoverCB(&ObjCoords))
	}
}

func ProjKeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	wWidth, wHeight := window.GetSize()

	if key == glfw.KeyBackspace && len(NameInputEnteredTxt) != 0 {
		NameInputEnteredTxt = NameInputEnteredTxt[:len(NameInputEnteredTxt)-1]
	} else if key == glfw.KeySpace {
		NameInputEnteredTxt += " "
	} else if key == glfw.KeyEnter && len(NameInputEnteredTxt) != 0 {
		// create file
		rootPath, _ := GetRootPath()

		ProjectName = NameInputEnteredTxt + ".s3p"
		outPath := filepath.Join(rootPath, ProjectName)
		os.WriteFile(outPath, []byte(""), 0777)

		// initialize
		SlideFormat = make([][]Drawn, 0)
		SlideFormat = append(SlideFormat, make([]Drawn, 0)) // first slide

		// move to work view
		DrawWorkView(window, CurrentSlide)
		window.SetMouseButtonCallback(workViewMouseCallback)
		window.SetKeyCallback(nil)
		window.SetScrollCallback(workViewScrollCB)
		return
	} else {
		NameInputEnteredTxt += glfw.GetKeyName(key, scancode)
	}

	nIRS := ProjObjCoords[PROJ_NameInput]
	theCtx := Continue2dCtx(CurrentWindowFrame, &ProjObjCoords)
	theCtx.drawInputB(PROJ_NameInput, nIRS.OriginX, nIRS.OriginY, nIRS.Width, NameInputEnteredTxt, true)

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, theCtx.ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	CurrentWindowFrame = theCtx.ggCtx.Image()
}
