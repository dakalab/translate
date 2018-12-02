package main

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

var errmsg string

func TestProcess(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tmpOutput := "/tmp/translate-result.json"
	os.Args = []string{"translate", "-i", "./demo.json", "-o", tmpOutput, "-t", "fr"}

	process()

	assert.Equal(t, "./demo.json", args.inputFile)
	assert.Equal(t, "fr", args.targetLang)

	_, err := os.Stat(tmpOutput)
	assert.NoError(t, err)
	os.Remove(tmpOutput)
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
	c.printSupportedLang()
	assert.Empty(t, errmsg)

	os.Setenv("GCLOUD_API_KEY", "fake-api-key")
	c = client(ctx)
	c.printSupportedLang()
	assert.NotEmpty(t, errmsg)
	errmsg = ""

	// restore the correct key
	os.Setenv("GCLOUD_API_KEY", key)
	c.Close()
}

func TestTranslate(t *testing.T) {
	key := os.Getenv("GCLOUD_API_KEY")
	ctx := context.Background()
	exit = func(err error) { errmsg = err.Error() }

	jsonMap := map[string]string{
		"greeting": "Hello world!",
	}

	c := client(ctx)
	c.translate(jsonMap, language.English, language.French)
	assert.Empty(t, errmsg)
	assert.Equal(t, "Bonjour le monde!", jsonMap["greeting"])

	os.Setenv("GCLOUD_API_KEY", "fake-api-key")
	c = client(ctx)
	c.translate(jsonMap, language.English, language.French)
	assert.NotEmpty(t, errmsg)
	errmsg = ""

	// restore the correct key
	os.Setenv("GCLOUD_API_KEY", key)
	c.Close()
}

func TestGetSourceJSON(t *testing.T) {
	exit = func(err error) { errmsg = err.Error() }

	res := getSourceJSON("./demo.json")
	assert.Equal(t, "Hello world!", res["greeting"])

	res = getSourceJSON("./not-exists.json")
	assert.NotEmpty(t, errmsg)
	errmsg = ""

	res = getSourceJSON("./invalid.json")
	assert.NotEmpty(t, errmsg)
	errmsg = ""
}

func TestSaveResult(t *testing.T) {
	exit = func(err error) { errmsg = err.Error() }

	jsonMap := map[string]string{
		"a": "a",
	}

	file := "/tmp/translate-result.json"
	saveResult(jsonMap, file)
	_, err := os.Stat(file)
	assert.NoError(t, err)
	os.Remove(file)

	file = "/tmp/readonly.json"
	ioutil.WriteFile(file, nil, 0444)
	saveResult(jsonMap, file)
	assert.NotEmpty(t, errmsg)
	os.Remove(file)
	errmsg = ""
}
