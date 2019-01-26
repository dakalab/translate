package translater

import (
	"context"
	"errors"
	"io/ioutil"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// HTMLTranslater - html translater
type HTMLTranslater struct {
	Client      *translate.Client
	Original    string
	Translation string
}

// NewHTMLTranslater - return a new HTMLTranslater
func NewHTMLTranslater(client *translate.Client) *HTMLTranslater {
	return &HTMLTranslater{Client: client}
}

// ParseFile - parse json file
func (trans *HTMLTranslater) ParseFile(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	trans.Original = string(content)
	return nil
}

// Translate - translate the original into translation
func (trans *HTMLTranslater) Translate(sl, tl language.Tag) error {
	if trans.Client == nil {
		return errors.New("translate client is nil")
	}

	var original = []string{
		trans.Original,
	}

	translations, err := trans.Client.Translate(context.Background(), original, tl, &translate.Options{
		Source: sl,
		Format: translate.HTML,
	})
	if err != nil {
		return err
	}

	trans.Translation = translations[0].Text

	return nil
}

// SaveResult - save translation to output file
func (trans *HTMLTranslater) SaveResult(file string) error {
	return ioutil.WriteFile(file, []byte(trans.Translation), 0644)
}
