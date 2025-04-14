package os

type Clipboard interface {
	SetText(text string) error
	GetText() (string, error)
}
