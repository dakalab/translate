package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type exiter func(err error)

type translateClient struct {
	*translate.Client
	ctx context.Context
}

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
		client.printSupportedLang()
		return
	}

	sl := language.MustParse(args.sourceLang)
	tl := language.MustParse(args.targetLang)

	jsonMap := getSourceJSON(args.inputFile)

	client.translate(jsonMap, sl, tl)

	saveResult(jsonMap, args.outputFile)

	if args.outputFile != "/dev/stdout" {
		log.Println("Translate successfully and save into " + args.outputFile)
	}
}

func parse() {
	flag.StringVar(&args.inputFile, "i", "", "the input path of json file to be translated")
	flag.StringVar(&args.outputFile, "o", "/dev/stdout", "the output path to save translated json file")
	flag.StringVar(&args.sourceLang, "s", "en", "source language")
	flag.StringVar(&args.targetLang, "t", "en", "target language")
	flag.BoolVar(&args.listLang, "l", false, "list available languages")

	flag.Parse()
}

func client(ctx context.Context) *translateClient {
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

	return &translateClient{client, ctx}
}

func getSourceJSON(file string) map[string]string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		exit(err)
		return nil
	}

	jsonMap := make(map[string]string)
	if err := json.Unmarshal(content, &jsonMap); err != nil {
		exit(err)
		return nil
	}

	return jsonMap
}

func saveResult(jsonMap map[string]string, file string) {
	result, err := json.MarshalIndent(jsonMap, "", "    ")
	if err != nil {
		exit(err)
		return
	}

	if err := ioutil.WriteFile(file, result, 0644); err != nil {
		exit(err)
		return
	}
}

func (client *translateClient) printSupportedLang() {
	langs, err := client.SupportedLanguages(client.ctx, language.English)
	if err != nil {
		exit(err)
		return
	}
	for _, lang := range langs {
		fmt.Println(lang.Tag.String() + ": " + lang.Name)
	}
	return
}

func (client *translateClient) translate(jsonMap map[string]string, sl, tl language.Tag) {
	var keys []string
	var values []string
	for k, v := range jsonMap {
		keys = append(keys, k)
		values = append(values, v)
	}
	translations, err := client.Translate(client.ctx, values, tl, &translate.Options{
		Source: sl,
	})
	if err != nil {
		exit(err)
		return
	}

	for i := 0; i < len(translations); i++ {
		jsonMap[keys[i]] = translations[i].Text
	}
}
