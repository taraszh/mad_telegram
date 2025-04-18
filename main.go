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
	oslocal "wacky_message/utils/os"
	utils "wacky_message/utils/os/windows"
	"wacky_message/window"

	"golang.design/x/hotkey"
)

const maxInputLength = 600
const minInputLength = 4
const triggerSuffix = "!!1"

//go:embed tray/systray_adapter/icon.ico
var iconData []byte

var processing sync.Mutex

var selectedWindows []string

func main() {
	if !isWindows() {
		println("At the moment, this application is supported exclusively on Windows.")
		return
	}

	selectedWindows, _ = oslocal.LoadSelectedWindows()

	var sysTray = &systray_adapter.SystrayAdapter{}
	var windowsUtils, _ = utils.NewWindow()
	var keyboard = utils.NewKeyboard()
	var hk = hotkey.New(nil, hotkey.KeyF2)

	var emojifier = emoji.NewEndStringEmojifier()
	var translator = translate.NewGoogleTranslator()
	var clipboard = utils.NewClipboard()

	messageQueue := make(chan string, 10)

	go sysTray.Run(
		func() { onReady(sysTray, windowsUtils) },
		func() { os.Exit(0) },
	)

	go processMessagesWithTrigger(clipboard, messageQueue)
	go processMessageByHotkey(hk, clipboard, keyboard, windowsUtils, messageQueue)

	processMessageQueue(messageQueue, keyboard, translator, emojifier)

	select {}
}

func onReady(trayAdapter tray.Tray, windowsUtil *utils.Window) {
	trayAdapter.SetIcon(iconData)
	trayAdapter.SetTooltip("wacky_message")
	trayAdapter.AddMenu("Quit", "", func() { trayAdapter.Quit() })
	trayAdapter.AddMenu(
		"Choose windows where hotkey will be active",
		"",
		func() {
			selectedWindows = window.NewInputWindow().
				OpenHotKeySettings(windowsUtil.WindowClassMap(), selectedWindows)

			if err := oslocal.SaveSelectedWindows(selectedWindows); err != nil {
				println("Error saving selected windows:", err)
			} else {
				println("Selected windows saved successfully.")
			}
		},
	)
}

func processMessagesWithTrigger(
	clipboard *utils.Clipboard,
	queue chan<- string,
) {
	fmt.Println("Clipboard trigger is registered. Waiting for messages...")

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
	clipboard *utils.Clipboard,
	keyboard *utils.Keyboard,
	windowsUtil *utils.Window,
	queue chan<- string,
) {
	if err := hotkey.Register(); err != nil {
		panic("failed to register F2: " + err.Error())
	}

	fmt.Println("Hotkey registered, waiting for F2...")

	for {
		<-hotkey.Keydown()
		fmt.Println("F2 pressed — action triggered!")

		if skipWindow(windowsUtil) {
			println("Skipping action due to window mismatch.")
			continue
		}

		keyboard.SendCtrlPlusKey(utils.VK_A)
		keyboard.SendCtrlPlusKey(utils.VK_X)

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

func skipWindow(windowsUtil *utils.Window) bool {
	skip := true

	for _, selectedWindow := range selectedWindows {
		if selectedWindow == windowsUtil.ForegroundWindowClass() {
			println("Skipping window:", selectedWindow)

			skip = false
			break
		}
	}

	return skip
}

func processMessageQueue(
	queue <-chan string,
	keyboard *utils.Keyboard,
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

func getMessageFromClipBoard(clipboard oslocal.Clipboard) string {
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

func clearClipboard(clipboard *utils.Clipboard) {
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
