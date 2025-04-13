package translate

import (
	"fmt"
	translate "github.com/gilang-as/google-translate"
)

type GoogleTranslatorStronger struct {
}

func NewGoogleTranslatorStronger() *GoogleTranslatorStronger {
	return &GoogleTranslatorStronger{}
}

func (g *GoogleTranslatorStronger) Translate(message string) (string, error) {
	fmt.Println("Translating message via google-translate")

	request := translate.Translate{
		Text: message,
		From: "uk",
		To:   "pl",
	}

	response, err := translate.Translator(request)
	if err != nil {
		return "", err
	}

	fmt.Println("Pronunciation :", response.Pronunciation)

	return response.Text, err
}
