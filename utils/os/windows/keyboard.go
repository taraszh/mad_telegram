package windows

import (
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

type Keyboard struct {
	dll        *syscall.LazyDLL
	keybdEvent *syscall.LazyProc
	sendInput  *syscall.LazyProc
}

func NewKeyboard() *Keyboard {
	dll := syscall.NewLazyDLL("user32.dll")
	return &Keyboard{
		dll:        dll,
		keybdEvent: dll.NewProc("keybd_event"),
		sendInput:  dll.NewProc("SendInput"),
	}
}

type keyboardInput struct {
	Type uint32
	Ki   keybdInput
	_    [8]byte
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

func (k *Keyboard) TypeMessage(text string) error {
	const VK_ENTER = 0x0D
	const VK_SHIFT = 0x10

	for _, r := range text {
		if r == '\n' {
			k.SendKeyDown(VK_SHIFT)
			time.Sleep(15 * time.Millisecond)
			k.SendKeyDown(VK_ENTER)
			time.Sleep(15 * time.Millisecond)
			k.SendKeyUp(VK_ENTER)
			time.Sleep(15 * time.Millisecond)
			k.SendKeyUp(VK_SHIFT)
			continue
		}

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

			k.sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))

			input.Ki.Flags = KEYEVENTF_UNICODE | KEYEVENTF_KEYUP
			k.sendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))

			time.Sleep(25 * time.Millisecond)
		}
	}

	return nil
}

func (k *Keyboard) SendKeyDown(key uint8) {
	k.keybdEvent.Call(uintptr(key), 0, 0, 0)
}

func (k *Keyboard) SendKeyUp(key uint8) {
	k.keybdEvent.Call(uintptr(key), 0, 2, 0)
}

func (k *Keyboard) SendCtrlA() {
	const VK_CONTROL = 0x11
	const VK_A = 0x41

	k.SendKeyDown(VK_CONTROL)
	time.Sleep(15 * time.Millisecond)
	k.SendKeyDown(VK_A)
	time.Sleep(15 * time.Millisecond)
	k.SendKeyUp(VK_A)
	time.Sleep(15 * time.Millisecond)
	k.SendKeyUp(VK_CONTROL)
}
