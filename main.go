package main

import (
	_ "embed"
	"fmt"
	"github.com/taraszh/mad_telegram/input"
	"github.com/taraszh/mad_telegram/tray"
	"github.com/taraszh/mad_telegram/tray/systray_adapter"
	"github.com/taraszh/mad_telegram/windows"
	"math"
	"os"
	"sync"
	"time"
)

const maxInputLength = 500

//go:embed tray/systray_adapter/icon.ico
var iconData []byte

// var window *app.Window

var showWindow = make(chan struct{}, 1)

var wg = &sync.WaitGroup{}

func main() {
	var sysTray tray.Tray = &systray_adapter.SystrayAdapter{}

	go func() {
		sysTray.Run(
			func() { onReady(sysTray) },
			func() { os.Exit(0) },
		)
	}()

	go func() {
		var message string
		var inputWindow *input.Window = input.NewInputWindow()
		var user32 = windows.NewUser32()

		for range showWindow {
			wg.Add(1)

			println(fmt.Sprintf("Trigger received. InputWindow: %v", inputWindow == nil))

			if inputWindow.GetWindow() == nil {
				defer wg.Done()

				message = inputWindow.OpenInputWindow()
				modifiedMessage := processMessage(message)

				if !sendInTelegram(modifiedMessage, user32) {
					println("Telegram's input is not active, saving message to clipboard")
					_ = windows.SetClipboardText(modifiedMessage)
				}

				clearSubmittedText(message)
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

func clearSubmittedText(message string) {
	message = ""
}

func processMessage(message string) string {
	if len(message) == 0 {
		println("No text submitted")
		return ""
	}

	message = message[:int(math.Min(float64(len(message)), float64(maxInputLength)))]

	return message + " - " + "🗡️"
}

func sendInTelegram(modifiedMessage string, user32 *windows.User32) bool {
	telegramHWND := user32.GetTelegramHWND()

	if telegramHWND == 0 {
		fmt.Println("Telegram window not detected.")
		return false
	}

	// Check if input is active, print text, copy it and compare with printed
	user32.SetForegroundWindow(telegramHWND)
	user32.SendCtrlX()
	time.Sleep(500 * time.Millisecond)

	user32.SendString(":::")
	user32.SendCtrlA()
	user32.SendCtrlX()
	clipBoardText, _ := windows.GetClipboardText()

	if clipBoardText != ":::" {
		return false
	}

	fmt.Println("Telegram's input is active, sending message")

	user32.SendString(modifiedMessage)

	return true
}
