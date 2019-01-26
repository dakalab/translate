package translater

import (
	"context"
	"errors"
	"io/ioutil"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// YAMLTranslater - yaml translater
type YAMLTranslater struct {
	Client      *translate.Client
	Original    map[string]string
	Translation map[string]string
}

// NewYAMLTranslater - return a new YAMLTranslater
func NewYAMLTranslater(client *translate.Client) *YAMLTranslater {
	return &YAMLTranslater{Client: client}
}

// ParseFile - parse yaml file
func (trans *YAMLTranslater) ParseFile(file string) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, &trans.Original)
}

// Translate - translate the original into translation
func (trans *YAMLTranslater) Translate(sl, tl language.Tag) error {
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
func (trans *YAMLTranslater) SaveResult(file string) error {
	result, err := yaml.Marshal(trans.Translation)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, result, 0644)
	return err
}
