package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/translate"
	"github.com/dakalab/translate/translater"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type exiter func(err error)

type arguments struct {
	inputFile  string
	outputFile string
	sourceLang string
	targetLang string
	listLang   bool
}

var exit exiter = func(err error) { log.Fatal(err) }

var args arguments

func main() {
	process()
}

func process() {
	parse()

	ctx := context.Background()

	client := client(ctx)
	defer client.Close() // close the client when finished.

	if args.listLang {
		printSupportedLang(client)
		return
	}

	if err := doTranslate(client, args); err != nil {
		exit(err)
	}

	if args.outputFile != "/dev/stdout" {
		log.Println("Translate successfully and save into " + args.outputFile)
	}
}

func doTranslate(client *translate.Client, args arguments) error {
	sl := language.MustParse(args.sourceLang)
	tl := language.MustParse(args.targetLang)

	var jsonTranslater = translater.NewJSONTranslater(client)
	err := jsonTranslater.ParseFile(args.inputFile)
	if err == nil {
		jsonTranslater.Translate(sl, tl)
		jsonTranslater.SaveResult(args.outputFile)
		return nil
	}

	var yamlTranslater = translater.NewYAMLTranslater(client)
	err = yamlTranslater.ParseFile(args.inputFile)
	if err == nil {
		yamlTranslater.Translate(sl, tl)
		yamlTranslater.SaveResult(args.outputFile)
		return nil
	}

	var htmlTranslater = translater.NewHTMLTranslater(client)
	err = htmlTranslater.ParseFile(args.inputFile)
	if err != nil {
		return err
	}
	htmlTranslater.Translate(sl, tl)
	htmlTranslater.SaveResult(args.outputFile)

	return nil
}

func parse() {
	flag.StringVar(&args.inputFile, "i", "", "the input file to be translated")
	flag.StringVar(&args.outputFile, "o", "/dev/stdout", "the output path to save translated file")
	flag.StringVar(&args.sourceLang, "s", "en", "source language")
	flag.StringVar(&args.targetLang, "t", "en", "target language")
	flag.BoolVar(&args.listLang, "l", false, "list available languages")

	flag.Parse()
}

func client(ctx context.Context) *translate.Client {
	apiKey := os.Getenv("GCLOUD_API_KEY")
	if apiKey == "" {
		exit(errors.New("You need to set google cloud API key in GCLOUD_API_KEY environment variable"))
		return nil
	}

	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		exit(err)
		return nil
	}

	return client
}

func printSupportedLang(client *translate.Client) {
	langs, err := client.SupportedLanguages(context.Background(), language.English)
	if err != nil {
		exit(err)
		return
	}
	for _, lang := range langs {
		fmt.Println(lang.Tag.String() + ": " + lang.Name)
	}
	return
}
