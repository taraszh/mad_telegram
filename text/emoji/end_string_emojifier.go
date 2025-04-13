package emoji

type EndStringEmojifier struct {
}

func NewEndStringEmojifier() *EndStringEmojifier {
	return &EndStringEmojifier{}
}

func (e *EndStringEmojifier) Emojify(message string) (string, error) {
	return message + " 🗡️", nil
}
