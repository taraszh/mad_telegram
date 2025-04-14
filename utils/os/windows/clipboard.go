package windows

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	cfUnicodeText  = 13
	ghnd           = 0x0042 // GMEM_MOVEABLE | GMEM_ZEROINIT
	memoryMaxBytes = 1 << 20
)

type Clipboard struct {
	user32   *syscall.LazyDLL
	kernel32 *syscall.LazyDLL
	procs    clipboardProcs
}

type clipboardProcs struct {
	openClipboard    *syscall.LazyProc
	closeClipboard   *syscall.LazyProc
	emptyClipboard   *syscall.LazyProc
	setClipboardData *syscall.LazyProc
	getClipboardData *syscall.LazyProc
	globalAlloc      *syscall.LazyProc
	globalLock       *syscall.LazyProc
	globalUnlock     *syscall.LazyProc
}

func NewClipboard() *Clipboard {
	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	return &Clipboard{
		user32:   user32,
		kernel32: kernel32,
		procs: clipboardProcs{
			openClipboard:    user32.NewProc("OpenClipboard"),
			closeClipboard:   user32.NewProc("CloseClipboard"),
			emptyClipboard:   user32.NewProc("EmptyClipboard"),
			setClipboardData: user32.NewProc("SetClipboardData"),
			getClipboardData: user32.NewProc("GetClipboardData"),
			globalAlloc:      kernel32.NewProc("GlobalAlloc"),
			globalLock:       kernel32.NewProc("GlobalLock"),
			globalUnlock:     kernel32.NewProc("GlobalUnlock"),
		},
	}
}

func (c *Clipboard) SetText(text string) error {
	if r, _, _ := c.procs.openClipboard.Call(0); r == 0 {
		return fmt.Errorf("failed to open clipboard")
	}
	defer c.procs.closeClipboard.Call()

	c.procs.emptyClipboard.Call()

	data, err := syscall.UTF16FromString(text)
	if err != nil {
		return err
	}
	dataSize := len(data) * 2

	hMem, _, _ := c.procs.globalAlloc.Call(ghnd, uintptr(dataSize))
	if hMem == 0 {
		return fmt.Errorf("failed to allocate memory")
	}

	ptr, _, _ := c.procs.globalLock.Call(hMem)
	if ptr == 0 {
		return fmt.Errorf("failed to lock memory")
	}
	defer c.procs.globalUnlock.Call(hMem)

	copy((*[memoryMaxBytes]uint16)(unsafe.Pointer(ptr))[:], data)

	if r, _, _ := c.procs.setClipboardData.Call(cfUnicodeText, hMem); r == 0 {
		return fmt.Errorf("failed to set clipboard data")
	}

	return nil
}

func (c *Clipboard) GetText() (string, error) {
	if r, _, _ := c.procs.openClipboard.Call(0); r == 0 {
		return "", fmt.Errorf("failed to open clipboard")
	}
	defer c.procs.closeClipboard.Call()

	hMem, _, _ := c.procs.getClipboardData.Call(cfUnicodeText)
	if hMem == 0 {
		return "", fmt.Errorf("no text data in clipboard")
	}

	ptr, _, _ := c.procs.globalLock.Call(hMem)
	if ptr == 0 {
		return "", fmt.Errorf("failed to lock memory")
	}
	defer c.procs.globalUnlock.Call(hMem)

	text := syscall.UTF16ToString((*[memoryMaxBytes]uint16)(unsafe.Pointer(ptr))[:])
	return text, nil
}
