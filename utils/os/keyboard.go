package os

type Keyboard interface {
	TypeMessage(string) error
	SendCtrlPlusKey(key byte)
}
