package funtranslator

type translationResponse struct {
	Contents struct {
		Translated string `json:"translated"`
	} `json:"contents"`
}
