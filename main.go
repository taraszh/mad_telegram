package main

import (
	_ "embed"
	"fmt"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/taraszh/mad_telegram/tray"
	"github.com/taraszh/mad_telegram/tray/systray_adapter"
	"image/color"
	"math"
	"os"
	"sync"
)

const maxInputLength = 5

//go:embed tray/systray_adapter/icon.ico
var iconData []byte

var window *app.Window
var showWindow = make(chan struct{}, 1)

var inputEditor widget.Editor
var submitButton widget.Clickable
var submittedText string

var wg sync.WaitGroup

func main() {
	var sysTray tray.Tray = &systray_adapter.SystrayAdapter{}

	go func() {
		sysTray.Run(
			func() { onReady(sysTray) },
			func() { os.Exit(0) },
		)
	}()

	go func() {
		for range showWindow {
			println(fmt.Sprintf("Trigger received. Is window nil: %v", window == nil))

			if window == nil {
				wg.Add(1)
				openInputWindow()
				processMessage()
				clearSubmittedText()
			} else {
				println("Window is already open, ignoring trigger.")
			}
		}
	}()

	select {}
}

func onReady(trayAdapter tray.Tray) {
	trayAdapter.SetIcon(iconData)
	trayAdapter.SetTooltip("wacky_message")
	trayAdapter.AddMenu(
		"Open message modifier",
		"Open input window",
		func() { showWindow <- struct{}{} },
	)

	trayAdapter.AddSeparator()

	trayAdapter.AddMenu(
		"Quit",
		"Exit app",
		func() { trayAdapter.Quit() },
	)
}

func openInputWindow() {
	defer wg.Done()

	window = new(app.Window)
	window.Option(app.Size(350, 250))
	window.Option(app.Title("🥴wacky_message"))

	theme := material.NewTheme()

	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			window = nil
			return
		case app.FrameEvent:
			//maroon := color.NRGBA{R: 127, G: 0, B: 0, A: 255}
			white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
			black := color.NRGBA{R: 0, G: 0, B: 0, A: 255}

			gtx := app.NewContext(&ops, e)

			title := material.H5(theme, "☠️Crazy Message🤪")
			title.Color = black
			title.Alignment = text.Middle

			paint.Fill(gtx.Ops, white)
			title.Layout(gtx)

			if submitButton.Clicked(gtx) {
				submittedText = inputEditor.Text()
				inputEditor.SetText("")

				window.Perform(system.ActionClose)
			}

			layout.Center.Layout(
				gtx,
				func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Vertical,
					}.Layout(
						gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Editor(theme, &inputEditor, "Enter a message...").Layout(gtx)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Button(theme, &submitButton, "Modify and Insert").Layout(gtx)
						}),
					)
				},
			)

			e.Frame(gtx.Ops)
		}
	}
}

func clearSubmittedText() {
	submittedText = ""
}

func processMessage() {
	if len(submittedText) == 0 {
		println("No text submitted")
		return
	}

	submittedText = submittedText[:int(math.Min(float64(len(submittedText)), float64(maxInputLength)))]

	println(submittedText + " - " + "🗡️")
}
