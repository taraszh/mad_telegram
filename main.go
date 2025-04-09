package main

import (
	_ "embed"
	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget/material"
	"github.com/getlantern/systray"
	"image/color"
	"os"
	"sync"
)

//go:embed icon.ico
var iconData []byte

var window *app.Window
var showWindow = make(chan struct{})

var wg sync.WaitGroup // Для синхронізації горутин

func main() {
	go func() {
		systray.Run(onReady, func() {
			os.Exit(0)
		})
	}()

	go func() {
		for range showWindow {
			if window == nil {
				wg.Add(1)
				openInputWindow()
			}
		}
	}()

	select {} // keep main alive
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTooltip("mad_telegram")
	systray.AddMenuItem("Send Message", "Open input window").ClickedCh = showWindow
	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Exit app")

	go func() {
		<-quit.ClickedCh
		systray.Quit()
	}()
}

func openInputWindow() {
	defer wg.Done()

	window = new(app.Window)
	theme := material.NewTheme()

	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			window = nil
			return
		case app.FrameEvent:
			maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

			gtx := app.NewContext(&ops, e)

			title := material.H3(theme, "☠️Crazy Message🤪")
			title.Color = maroon
			title.Alignment = text.Middle

			paint.Fill(gtx.Ops, black)

			title.Layout(gtx)

			// Pass the drawing operations to the GPU.
			e.Frame(gtx.Ops)
		}
	}
}
