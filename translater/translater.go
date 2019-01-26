package translater

import "golang.org/x/text/language"

// Translater - interface of translater
type Translater interface {
	ParseFile(file string) error
	Translate(sl, tl language.Tag) error
	SaveResult(file string) error
}
