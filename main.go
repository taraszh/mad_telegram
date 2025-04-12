package main

import (
	_ "embed"
	"fmt"
	"github.com/taraszh/mad_telegram/input"
	"github.com/taraszh/mad_telegram/tray"
	"github.com/taraszh/mad_telegram/tray/systray_adapter"
	"math"
	"os"
	"sync"
)

const maxInputLength = 500

//go:embed tray/systray_adapter/icon.ico
var iconData []byte

// var window *app.Window

var message string
var showWindow = make(chan struct{}, 1)

var wg = &sync.WaitGroup{}

func main() {
	var sysTray tray.Tray = &systray_adapter.SystrayAdapter{}
	var inputWindow *input.Window = input.NewInputWindow()

	go func() {
		sysTray.Run(
			func() { onReady(sysTray) },
			func() { os.Exit(0) },
		)
	}()

	go func() {
		for range showWindow {
			wg.Add(1)

			println(fmt.Sprintf("Trigger received. Is inputWindow nil: %v", inputWindow == nil))

			if inputWindow.GetWindow() == nil {
				defer wg.Done()

				message = inputWindow.OpenInputWindow()
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
	trayAdapter.AddMenu("Quit", "Exit app", func() { trayAdapter.Quit() })
}

func clearSubmittedText() {
	message = ""
}

func processMessage() {
	if len(message) == 0 {
		println("No text submitted")
		return
	}

	message = message[:int(math.Min(float64(len(message)), float64(maxInputLength)))]

	println(message + " - " + "🗡️")
}
