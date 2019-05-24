package translater

import (
	"encoding/json"
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
	return json.Unmarshal(content, &trans.Original)
}

// Translate - translate the original into translation
func (trans *JSONTranslater) Translate(sl, tl language.Tag) error {
	var err error
	trans.Translation, err = convert(trans.Client, trans.Original, sl, tl)
	return err
}

// SaveResult - save translation to output file
func (trans *JSONTranslater) SaveResult(file string) error {
	result, err := json.MarshalIndent(trans.Translation, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, result, 0644)
}
