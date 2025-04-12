package windows

import (
	"fmt"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

type User32 struct {
	dll                 *syscall.LazyDLL
	enumWindows         *syscall.LazyProc
	getWindowTextW      *syscall.LazyProc
	getClassNameW       *syscall.LazyProc
	isWindowVisible     *syscall.LazyProc
	showWindow          *syscall.LazyProc
	keybd_event         *syscall.LazyProc
	openClipboard       *syscall.LazyProc
	setForegroundWindow *syscall.LazyProc
	sendInput           *syscall.LazyProc
}

func NewUser32() *User32 {
	dll := syscall.NewLazyDLL("user32.dll")
	return &User32{
		dll:                 dll,
		enumWindows:         dll.NewProc("EnumWindows"),
		getWindowTextW:      dll.NewProc("GetWindowTextW"),
		getClassNameW:       dll.NewProc("GetClassNameW"),
		isWindowVisible:     dll.NewProc("IsWindowVisible"),
		showWindow:          dll.NewProc("ShowWindow"),
		keybd_event:         dll.NewProc("keybd_event"),
		openClipboard:       dll.NewProc("OpenClipboard"),
		setForegroundWindow: dll.NewProc("SetForegroundWindow"),
		sendInput:           dll.NewProc("SendInput"),
	}
}

// EnumProc WindowEnumProc is the callback signature for EnumWindows.
type EnumProc func(hwnd uintptr, lParam uintptr) uintptr

func (u *User32) EnumWindows(callback EnumProc, lparam uintptr) bool {
	ret, _, _ := u.enumWindows.Call(syscall.NewCallback(callback), lparam)
	return ret != 0
}

func (u *User32) GetWindowText(hwnd uintptr, buf []uint16) int {
	ret, _, _ := u.getWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(ret)
}

func (u *User32) IsWindowVisible(hwnd uintptr) bool {
	ret, _, _ := u.isWindowVisible.Call(hwnd)
	return ret != 0
}

func (u *User32) GetClassName(hwnd uintptr, buf []uint16) int {
	ret, _, _ := u.getClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(ret)
}

func (u *User32) MaximizeWindow(hwnd uintptr) {
	u.showWindow.Call(hwnd, uintptr(3)) // 3 - Maximize window
}

func (u *User32) OpenClipboard(hwnd uintptr) bool {
	ret, _, _ := u.openClipboard.Call(hwnd)
	return ret != 0
}

func (u *User32) GetTelegramHWND() uintptr {
	var telegramHWND uintptr

	callback := EnumProc(
		func(hwnd uintptr, lParam uintptr) uintptr {
			if !u.IsWindowVisible(hwnd) {
				return 1
			}

			buf := make([]uint16, 256)
			length := u.GetWindowText(hwnd, buf)
			if length > 0 {
				title := syscall.UTF16ToString(buf[:length])

				classBuf := make([]uint16, 256)
				classNameLen := u.GetClassName(hwnd, classBuf)
				className := syscall.UTF16ToString(classBuf[:classNameLen])

				index := strings.Index(title, "@")
				if index != -1 {
					fmt.Println("Telegram window found!")
					fmt.Printf("Window Title: %s | Window Class: %s\n", title, className)

					telegramHWND = hwnd
				}
			}

			return 1
		},
	)

	// Перераховуємо всі вікна
	if !u.EnumWindows(callback, 0) {
		fmt.Println("Error enumerating windows.")
	}

	return telegramHWND
}

func (u *User32) SetForegroundWindow(hwnd uintptr) bool {
	println("Setting foreground window")
	ret, _, _ := u.setForegroundWindow.Call(hwnd)

	time.Sleep(500 * time.Millisecond)

	return ret != 0
}

type keyboardInput struct {
	Type uint32
	Ki   keybdInput
	_    [8]byte // padding for alignment
}

type keybdInput struct {
	Vk        uint16
	Scan      uint16
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

const (
	INPUT_KEYBOARD    = 1
	KEYEVENTF_UNICODE = 0x0004
	KEYEVENTF_KEYUP   = 0x0002
)

func (u *User32) SendString(text string) {
	for _, r := range text {
		// Convert the rune to UTF-16 encoding to handle surrogate pairs
		utf16Encoded := utf16.Encode([]rune{r})

		for _, codeUnit := range utf16Encoded {
			input := keyboardInput{
				Type: INPUT_KEYBOARD,
				Ki: keybdInput{
					Scan:      codeUnit,
					Flags:     KEYEVENTF_UNICODE,
					ExtraInfo: 0,
				},
			}

			// Send key down
			u.sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))

			// Send key up
			input.Ki.Flags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP
			u.sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))

			time.Sleep(25 * time.Millisecond)
		}
	}
}

func (u *User32) SendKeyDown(key uint8) {
	u.keybd_event.Call(uintptr(key), 0, 0, 0)
}

func (u *User32) SendKeyUp(key uint8) {
	u.keybd_event.Call(uintptr(key), 0, 2, 0)
}

func (u *User32) SendCtrlA() {
	const VK_CONTROL = 0x11
	const VK_A = 0x41

	u.SendKeyDown(VK_CONTROL)
	time.Sleep(15 * time.Millisecond)
	u.SendKeyDown(VK_A)
	time.Sleep(15 * time.Millisecond)
	u.SendKeyUp(VK_A)
	time.Sleep(15 * time.Millisecond)
	u.SendKeyUp(VK_CONTROL)
}

func (u *User32) SendCtrlX() {
	const VK_CONTROL = 0x11
	const VK_X = 0x58

	u.SendKeyDown(VK_CONTROL)
	time.Sleep(15 * time.Millisecond)
	u.SendKeyDown(VK_X)
	time.Sleep(15 * time.Millisecond)
	u.SendKeyUp(VK_X)
	time.Sleep(15 * time.Millisecond)
	u.SendKeyUp(VK_CONTROL)
}
