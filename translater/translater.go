package translater

import (
	"context"
	"errors"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

// Translater - interface of translater
type Translater interface {
	ParseFile(file string) error
	Translate(sl, tl language.Tag) error
	SaveResult(file string) error
}

// max number of segments per request
const maxSegmentsPerRequest = 100

func convert(client *translate.Client, original map[string]string, sl, tl language.Tag) (map[string]string, error) {
	if client == nil {
		return nil, errors.New("translate client is nil")
	}

	var keys []string
	var values []string
	for k, v := range original {
		keys = append(keys, k)
		values = append(values, v)
	}

	var total = len(values)
	var translations []translate.Translation

	// there is limitation of the number of text segments, so we process maxSegmentsPerRequest segments each time
	for i := 0; i < total; i += maxSegmentsPerRequest {
		j := i + maxSegmentsPerRequest
		if j > total {
			j = total
		}
		parts, err := client.Translate(context.Background(), values[i:j], tl, &translate.Options{
			Source: sl,
		})
		if err != nil {
			return nil, err
		}
		translations = append(translations, parts...)
	}

	var result = make(map[string]string, len(translations))
	for i := 0; i < len(translations); i++ {
		result[keys[i]] = translations[i].Text
	}

	return result, nil
}
