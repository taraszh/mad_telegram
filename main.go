package main

import (
	_ "embed"
	"math"
	"os"
	"strings"
	"time"
	"wacky_message/tray"
	"wacky_message/tray/systray_adapter"
	"wacky_message/windows"
)

const maxInputLength = 1000
const minInputLength = 6
const triggerTemplate = "!!1"

//go:embed tray/systray_adapter/icon.ico
var iconData []byte

func main() {
	var sysTray = &systray_adapter.SystrayAdapter{}
	var user32 = windows.NewUser32()
	var message string

	go func() {
		sysTray.Run(
			func() { onReady(sysTray) },
			func() { os.Exit(0) },
		)
	}()

	go func() {
		for {
			message = getMessageFromClipBoard()

			if message != "" {
				user32.SendString(modifyMessage(message))
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {}
}

func onReady(trayAdapter tray.Tray) {
	trayAdapter.SetIcon(iconData)
	trayAdapter.SetTooltip("wacky_message")
	trayAdapter.AddMenu("Quit", "", func() { trayAdapter.Quit() })
}

func getMessageFromClipBoard() string {
	message, _ := windows.GetClipboardText()
	clean := strings.ReplaceAll(message, "\t", "")
	clean = strings.ReplaceAll(clean, "\n", "")
	clean = strings.ReplaceAll(clean, "\r", "")

	if strings.HasSuffix(clean, triggerTemplate) && len(clean) > minInputLength {
		println("Message found by template")

		err := windows.SetClipboardText("Modifying message 🗡️")

		if err != nil {
			println("Error setting clipboard text:", err)
		}

		return clean
	}

	return ""
}

func modifyMessage(message string) string {
	println("Modifying message")

	message = strings.TrimSuffix(message, triggerTemplate)
	message = message[:int(math.Min(float64(len(message)), float64(maxInputLength)))]

	modifiedMessage := message + " - " + "🗡️"

	println("Modified message: ", modifiedMessage)

	return modifiedMessage
}
