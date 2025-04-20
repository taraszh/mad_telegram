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
	osUtils "wacky_message/utils/os"
	clipboardUtils "wacky_message/utils/os/windows/clipboard"
	keyboardUtils "wacky_message/utils/os/windows/keyboard"
	windowUtils "wacky_message/utils/os/windows/window"
	"wacky_message/window"

	hk "golang.design/x/hotkey"
)

const maxInputLength = 600
const minInputLength = 4
const triggerSuffix = "!!1"

//go:embed tray/systray_adapter/icon.ico
var iconData []byte

var sysTray tray.Tray
var logWindow window.Log
var windowsUtils osUtils.Window
var keyboard osUtils.Keyboard
var clipboard osUtils.Clipboard
var hotkey *hk.Hotkey
var emojifier *emoji.EndStringEmojifier
var translator *translate.GoogleTranslator
var processing sync.Mutex
var selectedWindows []string

func main() {
	ensureWindowsCompatibility()
	initialize()

	messageQueue := make(chan string, 10)

	go sysTray.Run(
		func() { onReady(sysTray, windowsUtils) },
		func() { os.Exit(0) },
	)

	go processMessagesWithTrigger(messageQueue)
	go processMessageByHotkey(messageQueue)

	processMessageQueue(messageQueue)

	select {}
}

func ensureWindowsCompatibility() {
	if runtime.GOOS == "windows" {
		return
	}

	log("At the moment, this application is supported exclusively on Windows.")
	os.Exit(0)
}

func initialize() {
	logWindow = *window.NewLog()
	sysTray = systray_adapter.NewSystrayAdapter()

	windowsUtils = windowUtils.NewWindow()
	keyboard = keyboardUtils.NewKeyboard()
	clipboard = clipboardUtils.NewClipboard()

	hotkey = hk.New(nil, hk.KeyF2)

	emojifier = emoji.NewEndStringEmojifier()
	translator = translate.NewGoogleTranslator()

	selectedWindows, _ = osUtils.LoadSelectedWindows()
}

func onReady(trayAdapter tray.Tray, windowsUtil osUtils.Window) {
	trayAdapter.SetIcon(iconData)
	trayAdapter.SetTooltip("wacky_message")
	trayAdapter.AddMenu("Quit", "", func() { trayAdapter.Quit() })
	trayAdapter.AddSeparator()
	trayAdapter.AddMenu(
		"Log",
		"",
		func() {
			logWindow.Show()
		},
	)
	trayAdapter.AddMenu(
		"Hotkey settings",
		"",
		func() {
			configureHotkeySettings(windowsUtil)
		},
	)
}

func configureHotkeySettings(windowsUtil osUtils.Window) {
	selectedWindows = window.NewWindow().OpenHotKeySettings(windowsUtil.WindowClassMap(), selectedWindows)

	if err := osUtils.SaveSelectedWindows(selectedWindows); err != nil {
		log("Error saving selected windows:" + err.Error())
	} else {
		log("Selected windows saved successfully.")
	}
}

func processMessagesWithTrigger(queue chan<- string) {
	log("Clipboard trigger is registered. Waiting for messages...")

	if err := clipboard.SetText(""); err != nil {
		log("Error setting clipboard text: " + err.Error())
	}

	for {
		message := getMessageFromClipBoard(clipboard)

		if isTriggerPresent(message) {
			log("Message found in clipboard by template.")
			spamClipboard(clipboard)

			processing.Lock()
			queue <- message
			processing.Unlock()
		}

		time.Sleep(1 * time.Second)
	}
}

func processMessageByHotkey(queue chan<- string) {
	if err := hotkey.Register(); err != nil {
		log("Error registering hotkey: " + err.Error())
		return
	}

	log("Hotkey registered, waiting for F2...")

	for {
		<-hotkey.Keydown()
		log("F2 pressed — action triggered!")

		if skipWindow() {
			log("Skipping message processing due to window mismatch.")
			continue
		}

		keyboard.SendCtrlPlusKey(keyboardUtils.VK_A)
		keyboard.SendCtrlPlusKey(keyboardUtils.VK_X)

		time.Sleep(20 * time.Millisecond)

		message := getMessageFromClipBoard(clipboard)

		if isMessageValid(message) {
			log("Message found in clipboard.")

			processing.Lock()
			queue <- message
			processing.Unlock()
		}
	}
}

func skipWindow() bool {
	skip := true

	for _, selectedWindow := range selectedWindows {
		if selectedWindow == windowsUtils.ForegroundWindowClass() {
			log("Skipping window: " + selectedWindow)

			skip = false
			break
		}
	}

	return skip
}

func processMessageQueue(queue <-chan string) {
	for message := range queue {
		modified := modifyMessage(message)

		if err := keyboard.TypeMessage(modified); err != nil {
			log("Error typing message: " + err.Error())
		}
	}
}

func getMessageFromClipBoard(clipboard osUtils.Clipboard) string {
	message, err := clipboard.GetText()

	if err != nil {
		log(err.Error())
	}

	clean := strings.ReplaceAll(message, "\t", "")
	clean = strings.ReplaceAll(clean, "\r", "")

	return clean
}

func modifyMessage(originalMessage string) string {
	originalMessage = strings.TrimSuffix(originalMessage, triggerSuffix)
	originalMessage = originalMessage[:int(math.Min(float64(len(originalMessage)), float64(maxInputLength)))]

	log("Original message: " + originalMessage)

	translatedMessage, _ := translator.Translate(originalMessage)
	translatedEmojifiedMessage, _ := emojifier.Emojify(translatedMessage)
	originalEmojifiedMessage, _ := emojifier.Emojify(originalMessage)

	output := originalEmojifiedMessage + "\n" + translatedEmojifiedMessage

	log("Modified message: " + translatedEmojifiedMessage)

	return output
}

func spamClipboard(clipboard osUtils.Clipboard) {
	err := clipboard.SetText("Knife in your clipboard 🗡️")

	if err != nil {
		log("Error setting clipboard text: " + err.Error())
	}
}

func isTriggerPresent(message string) bool {
	if strings.HasSuffix(message, triggerSuffix) && isMessageValid(message) {
		return true
	}

	return false
}

func isMessageValid(message string) bool {
	return len(message) > minInputLength
}

func log(str string) {
	fmt.Println(str)
	logWindow.AddLog(str)
}
