package windows

import (
	"fmt"
	"syscall"
	"unsafe"
)

func SetClipboardText(text string) error {
	const CF_UNICODETEXT = 13

	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	openClipboard := user32.NewProc("OpenClipboard")
	closeClipboard := user32.NewProc("CloseClipboard")
	emptyClipboard := user32.NewProc("EmptyClipboard")
	setClipboardData := user32.NewProc("SetClipboardData")
	globalAlloc := kernel32.NewProc("GlobalAlloc")
	globalLock := kernel32.NewProc("GlobalLock")
	globalUnlock := kernel32.NewProc("GlobalUnlock")

	// Open clipboard
	r, _, _ := openClipboard.Call(0)
	if r == 0 {
		return fmt.Errorf("failed to open clipboard")
	}
	defer closeClipboard.Call()

	// Empty current clipboard
	emptyClipboard.Call()

	// Convert string to UTF-16
	data, err := syscall.UTF16FromString(text)
	if err != nil {
		return err
	}
	dataSize := len(data) * 2 // size in bytes

	// Allocate global memory
	hMem, _, _ := globalAlloc.Call(0x0042, uintptr(dataSize))
	if hMem == 0 {
		return fmt.Errorf("failed to allocate memory")
	}

	// Lock the memory to get a pointer
	ptr, _, _ := globalLock.Call(hMem)
	if ptr == 0 {
		return fmt.Errorf("failed to lock memory")
	}
	defer globalUnlock.Call(hMem)

	// Copy data to the memory
	copy((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:], data)

	// Set data to clipboard
	r, _, _ = setClipboardData.Call(CF_UNICODETEXT, hMem)
	if r == 0 {
		return fmt.Errorf("failed to set clipboard data")
	}

	return nil
}

func GetClipboardText() (string, error) {
	const CF_UNICODETEXT = 13

	if ok, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("OpenClipboard").Call(0); ok == 0 {
		return "", fmt.Errorf("failed to open clipboard")
	}
	defer syscall.NewLazyDLL("user32.dll").NewProc("CloseClipboard").Call()

	hMem, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("GetClipboardData").Call(CF_UNICODETEXT)
	if hMem == 0 {
		return "", fmt.Errorf("no text data in clipboard")
	}

	ptr, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalLock").Call(hMem)
	if ptr == 0 {
		return "", fmt.Errorf("failed to lock global memory")
	}
	defer syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalUnlock").Call(hMem)

	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(ptr))[:])

	return text, nil
}
