package os

type Window interface {
	WindowClassMap() map[string]string
	ForegroundWindowClass() string
}
