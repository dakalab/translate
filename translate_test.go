package main

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errmsg string

func TestProcess(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tmpOutput := "/tmp/translate-result.json"
	os.Args = []string{"translate", "-i", "./testfiles/demo.json", "-o", tmpOutput, "-t", "fr"}

	process()

	assert.Equal(t, "./testfiles/demo.json", args.inputFile)
	assert.Equal(t, "fr", args.targetLang)

	_, err := os.Stat(tmpOutput)
	assert.NoError(t, err)
	os.Remove(tmpOutput)
}

func TestYAML(t *testing.T) {
	tmpOutput := "/tmp/translate-result.yml"
	var args = arguments{
		inputFile:  "./testfiles/demo.yml",
		outputFile: tmpOutput,
		sourceLang: "en",
		targetLang: "fr",
	}

	err := doTranslate(client(context.Background()), args)
	assert.NoError(t, err)

	_, err = os.Stat(tmpOutput)
	assert.NoError(t, err)
	os.Remove(tmpOutput)
}

func TestHTML(t *testing.T) {
	tmpOutput := "/tmp/translate-result.html"
	var args = arguments{
		inputFile:  "./testfiles/demo.html",
		outputFile: tmpOutput,
		sourceLang: "en",
		targetLang: "fr",
	}

	err := doTranslate(client(context.Background()), args)
	assert.NoError(t, err)

	_, err = os.Stat(tmpOutput)
	assert.NoError(t, err)
	os.Remove(tmpOutput)

	// should fail with non-existent file
	args.inputFile = "./not-exists.html"
	err = doTranslate(client(context.Background()), args)
	assert.Error(t, err)
}

func TestClient(t *testing.T) {
	key := os.Getenv("GCLOUD_API_KEY")
	ctx := context.Background()
	exit = func(err error) { errmsg = err.Error() }

	os.Setenv("GCLOUD_API_KEY", "")
	client(ctx)
	assert.Equal(t, "You need to set google cloud API key in GCLOUD_API_KEY environment variable", errmsg)
	errmsg = ""

	// restore the correct key
	os.Setenv("GCLOUD_API_KEY", key)
	c := client(ctx)
	c.Close()
}

func TestPrintSupportedLang(t *testing.T) {
	key := os.Getenv("GCLOUD_API_KEY")
	ctx := context.Background()
	exit = func(err error) { errmsg = err.Error() }

	c := client(ctx)
	printSupportedLang(c)
	assert.Empty(t, errmsg)

	os.Setenv("GCLOUD_API_KEY", "fake-api-key")
	c = client(ctx)
	printSupportedLang(c)
	assert.NotEmpty(t, errmsg)
	errmsg = ""

	// restore the correct key
	os.Setenv("GCLOUD_API_KEY", key)
	c.Close()
}
