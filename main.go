package main

import (
	_ "embed"
	windows_go "golang.org/x/sys/windows"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
	"wacky_message/text/emoji"
	"wacky_message/text/translate"
	"wacky_message/tray"
	"wacky_message/tray/systray_adapter"
	os_local "wacky_message/utils/os"
	"wacky_message/utils/os/windows"
)

const maxInputLength = 1000
const minInputLength = 6
const triggerSuffix = "!!1"

//go:embed tray/systray_adapter/icon.ico
var iconData []byte

func main() {
	if !isWindows() {
		println("At the moment, this application is supported exclusively on Windows.")
		return
	}

	var sysTray = &systray_adapter.SystrayAdapter{}
	var keyboard = windows.NewKeyboard()

	var emojifier = emoji.NewEndStringEmojifier()
	var translator = translate.NewGoogleTranslator()

	var clipboard = windows.NewClipboard()

	go func() {
		sysTray.Run(
			func() { onReady(sysTray) },
			func() { os.Exit(0) },
		)
	}()

	processMessagesWithTrigger(clipboard, keyboard, translator, emojifier)

	select {}
}

func processMessagesWithTrigger(clipboard *windows.Clipboard, keyboard *windows.Keyboard, translator *translate.GoogleTranslator, emojifier *emoji.EndStringEmojifier) {
	var message string

	go func() {
		for {
			message = getMessageFromClipBoard(clipboard)

			if message != "" {
				println("Message found in clipboard")

				err := keyboard.TypeMessage(modifyMessage(message, translator, emojifier))
				if err != nil {
					println("Error typing message:", err)
				}

				fg_window := windows_go.GetForegroundWindow()
				if fg_window != 0 {
					println("FG window found")
				}
			}

			time.Sleep(200 * time.Millisecond)
		}
	}()
}

func onReady(trayAdapter tray.Tray) {
	trayAdapter.SetIcon(iconData)
	trayAdapter.SetTooltip("wacky_message")
	trayAdapter.AddMenu("Quit", "", func() { trayAdapter.Quit() })
}

func getMessageFromClipBoard(clipboard os_local.Clipboard) string {
	message, _ := clipboard.GetText()
	clean := strings.ReplaceAll(message, "\t", "")
	clean = strings.ReplaceAll(clean, "\r", "")

	if strings.HasSuffix(clean, triggerSuffix) && len(clean) > minInputLength {
		println("Message found by template")

		err := clipboard.SetText("Modifying message 🗡️")

		if err != nil {
			println("Error setting clipboard text:", err)
		}

		return clean
	}

	return ""
}

func modifyMessage(
	originalMessage string,
	translator translate.Translator,
	emojifier emoji.Emojifier,
) string {
	println("Modifying message")

	originalMessage = strings.TrimSuffix(originalMessage, triggerSuffix)
	originalMessage = originalMessage[:int(math.Min(float64(len(originalMessage)), float64(maxInputLength)))]

	println("Original message: ", originalMessage)

	translatedMessage, _ := translator.Translate(originalMessage)
	translatedEmojifiedMessage, _ := emojifier.Emojify(translatedMessage)
	originalEmojifiedMessage, _ := emojifier.Emojify(originalMessage)

	output := originalEmojifiedMessage + "\n" + translatedEmojifiedMessage

	println("Modified message: ", translatedEmojifiedMessage)

	return output
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}
