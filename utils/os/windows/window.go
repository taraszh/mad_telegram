package windows

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

type Window struct {
	user32              *syscall.LazyDLL
	isWindowVisible     *syscall.LazyProc
	getWindowTextW      *syscall.LazyProc
	getClassNameW       *syscall.LazyProc
	enumWindows         *syscall.LazyProc
	setForegroundWindow *syscall.LazyProc
	showWindow          *syscall.LazyProc
}

func NewWindow() (*Window, error) {
	user32 := syscall.NewLazyDLL("user32.dll")

	return &Window{
		user32:              user32,
		isWindowVisible:     user32.NewProc("IsWindowVisible"),
		getWindowTextW:      user32.NewProc("GetWindowTextW"),
		getClassNameW:       user32.NewProc("GetClassNameW"),
		enumWindows:         user32.NewProc("EnumWindows"),
		showWindow:          user32.NewProc("ShowWindow"),
		setForegroundWindow: user32.NewProc("SetForegroundWindow"),
	}, nil
}

func (w *Window) WindowHWND(substring string) uintptr {
	var window uintptr

	callback := EnumProc(
		func(hwnd uintptr, lParam uintptr) uintptr {
			if !w.IsWindowVisible(hwnd) {
				return 1
			}

			buf := make([]uint16, 256)
			length := w.GetWindowText(hwnd, buf)
			if length == 0 {
				return 1
			}

			title := syscall.UTF16ToString(buf[:length])

			className := w.ClassNameString(hwnd)

			if strings.Contains(strings.ToLower(title), substring) {
				fmt.Println("Window found!")
				fmt.Printf("Window Title: %s | Window Class: %s\n", title, className)
				window = hwnd
			}

			return 1
		},
	)

	if !w.EnumWindows(callback, 0) {
		fmt.Println("Error enumerating windows.")
	}

	return window
}

func (w Window) WindowClassMap() map[string]string {
	winClassMap := make(map[string]string)

	callback := EnumProc(
		func(hwnd uintptr, lParam uintptr) uintptr {
			if !w.IsWindowVisible(hwnd) {
				return 1
			}

			buf := make([]uint16, 256)
			length := w.GetWindowText(hwnd, buf)
			if length == 0 {
				return 1
			}

			title := syscall.UTF16ToString(buf[:length])
			className := w.ClassNameString(hwnd)

			winClassMap[title] = className

			//fmt.Printf("Window Title: %s | Window Class: %s\n", title, className)

			return 1
		},
	)

	if !w.EnumWindows(callback, 0) {
		fmt.Println("Error enumerating windows.")
	}

	return winClassMap
}

func (w *Window) ClassNameString(hwnd uintptr) string {
	buf := make([]uint16, 256)
	classNameLen, _, _ := w.getClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	className := syscall.UTF16ToString(buf[:classNameLen])

	return className
}

func (w *Window) IsWindowVisible(hwnd uintptr) bool {
	ret, _, _ := w.isWindowVisible.Call(hwnd)
	return ret != 0
}

func (w *Window) GetWindowText(hwnd uintptr, buf []uint16) int {
	ret, _, _ := w.getWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(ret)
}

func (w *Window) ClassName(hwnd uintptr, buf []uint16) int {
	ret, _, _ := w.getClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return int(ret)
}

type EnumProc func(hwnd uintptr, lParam uintptr) uintptr

func (w *Window) EnumWindows(callback EnumProc, lparam uintptr) bool {
	ret, _, _ := w.enumWindows.Call(syscall.NewCallback(callback), lparam)
	return ret != 0
}

func (w *Window) MaximizeWindow(hwnd uintptr) {
	_, _, _ = w.showWindow.Call(hwnd, uintptr(3)) // 3 - Maximize window
}

func (w *Window) SetForegroundWindow(hwnd uintptr) bool {
	println("Setting foreground window")
	ret, _, _ := w.setForegroundWindow.Call(hwnd)

	return ret != 0
}

func (w Window) ForegroundWindowClass() string {
	hwnd := win.GetForegroundWindow()
	if hwnd == 0 {
		return ""
	}

	return w.ClassNameString(uintptr(hwnd))
}
