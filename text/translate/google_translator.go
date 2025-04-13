package translate

import (
	"fmt"
	gtranslator "github.com/bas24/googletranslatefree"
)

type GoogleTranslator struct {
}

func NewGoogleTranslator() *GoogleTranslator {
	return &GoogleTranslator{}
}

func (g *GoogleTranslator) Translate(message string) (string, error) {
	fmt.Println("Translating message via googletranslatefree")

	translated, err := gtranslator.Translate(message, "uk", "pl")

	if err != nil {
		return err.Error(), err
	}

	return translated, err
}
