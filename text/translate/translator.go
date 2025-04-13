package translate

type Translator interface {
	Translate(message string) (string, error)
}
