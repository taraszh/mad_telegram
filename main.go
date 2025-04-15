package main

import (
	_ "embed"
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"wacky_message/text/emoji"
	"wacky_message/text/translate"
	"wacky_message/tray"
	"wacky_message/tray/systray_adapter"
	os_local "wacky_message/utils/os"
	"wacky_message/utils/os/windows"
	"golang.design/x/hotkey"
)

const maxInputLength = 300
const minInputLength = 4
const triggerSuffix = "!!1"

//go:embed tray/systray_adapter/icon.ico
var iconData []byte

var processing sync.Mutex

func main() {
	if !isWindows() {
		println("At the moment, this application is supported exclusively on Windows.")
		return
	}

	var sysTray = &systray_adapter.SystrayAdapter{}
	var keyboard = windows.NewKeyboard()
	var hk = hotkey.New(nil, hotkey.KeyF1)

	var emojifier = emoji.NewEndStringEmojifier()
	var translator = translate.NewGoogleTranslator()
	var clipboard = windows.NewClipboard()

	messageQueue := make(chan string, 10)

	go sysTray.Run(
		func() { onReady(sysTray) },
		func() { os.Exit(0) },
	)

	go processMessagesWithTrigger(clipboard, messageQueue)
	go processMessageByHotkey(hk, clipboard, keyboard, messageQueue)

	processMessageQueue(messageQueue, keyboard, translator, emojifier)

	select {}
}

func processMessagesWithTrigger(
	clipboard *windows.Clipboard,
	queue chan<- string,
) {
	_ = clipboard.SetText("")

	for {
		message := getMessageFromClipBoard(clipboard)

		if isTriggerPresent(message) && message != "" {
			processing.Lock()

			clearClipboard(clipboard)

			println("Message found in clipboard by template")
			queue <- message

			processing.Unlock()
		}

		time.Sleep(1 * time.Second)
	}
}

func processMessageByHotkey(
	hotkey *hotkey.Hotkey,
	clipboard *windows.Clipboard,
	keyboard *windows.Keyboard,
	queue chan<- string,
) {
	if err := hotkey.Register(); err != nil {
		panic("failed to register F1: " + err.Error())
	}

	fmt.Println("Hotkey F1 is active. Press it!")

	for {
		<-hotkey.Keydown()
		fmt.Println("F1 pressed — action triggered!")

		keyboard.SendCtrlPlusKey(windows.VK_A)
		keyboard.SendCtrlPlusKey(windows.VK_X)

		time.Sleep(20 * time.Millisecond)

		message := getMessageFromClipBoard(clipboard)

		if message != "" {
			processing.Lock()

			clearClipboard(clipboard)

			println("Message found in clipboard after hotkey")
			queue <- message

			processing.Unlock()
		}
	}
}

func processMessageQueue(
	queue <-chan string,
	keyboard *windows.Keyboard,
	translator *translate.GoogleTranslator,
	emojifier *emoji.EndStringEmojifier,
) {
	for message := range queue {
		modified := modifyMessage(message, translator, emojifier)
		err := keyboard.TypeMessage(modified)
		if err != nil {
			println("Error typing message:", err)
		}
	}
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

	return clean
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

func clearClipboard(clipboard *windows.Clipboard) {
	err := clipboard.SetText("Knife in your clipboard 🗡️")

	if err != nil {
		println("Error setting clipboard text:", err)
	}
}

func isTriggerPresent(message string) bool {
	if strings.HasSuffix(message, triggerSuffix) && len(message) > minInputLength {
		return true
	}

	return false
}
