package emoji

type Emojifier interface {
	Emojify(message string) (string, error)
}
