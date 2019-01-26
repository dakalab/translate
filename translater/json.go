package translater

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// JSONTranslater - json translater
type JSONTranslater struct {
	Client      *translate.Client
	Original    map[string]string
	Translation map[string]string
}

// NewJSONTranslater - return a new JSONTranslater
func NewJSONTranslater(client *translate.Client) *JSONTranslater {
	return &JSONTranslater{Client: client}
}

// ParseFile - parse json file
func (trans *JSONTranslater) ParseFile(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	trans.Original = make(map[string]string)
	err = json.Unmarshal(content, &trans.Original)
	return err
}

// Translate - translate the original into translation
func (trans *JSONTranslater) Translate(sl, tl language.Tag) error {
	if trans.Client == nil {
		return errors.New("translate client is nil")
	}

	var keys []string
	var values []string
	for k, v := range trans.Original {
		keys = append(keys, k)
		values = append(values, v)
	}
	translations, err := trans.Client.Translate(context.Background(), values, tl, &translate.Options{
		Source: sl,
	})
	if err != nil {
		return err
	}

	trans.Translation = make(map[string]string, len(translations))
	for i := 0; i < len(translations); i++ {
		trans.Translation[keys[i]] = translations[i].Text
	}

	return nil
}

// SaveResult - save translation to output file
func (trans *JSONTranslater) SaveResult(file string) error {
	result, err := json.MarshalIndent(trans.Translation, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, result, 0644)
	return err
}
