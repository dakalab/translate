package translater

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"cloud.google.com/go/translate"
	"github.com/stretchr/testify/suite"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type YAMLTestSuite struct {
	suite.Suite
	translater *YAMLTranslater
}

func TestYAMLTestSuite(t *testing.T) {
	suite.Run(t, new(YAMLTestSuite))
}

// SetupSuite - run by testify once at the very start of the testing suite, before any tests are run.
func (suite *YAMLTestSuite) SetupSuite() {
	client, err := translate.NewClient(
		context.TODO(),
		option.WithAPIKey(os.Getenv("GCLOUD_API_KEY")),
	)
	suite.NoError(err)
	suite.translater = NewYAMLTranslater(client)
}

// TearDownSuite - run by testify once, at the very end of the testing suite, after all tests have been run.
func (suite *YAMLTestSuite) TearDownSuite() {
	suite.translater.Client.Close()
}

// SetupTest - run before every test in the suite.
func (suite *YAMLTestSuite) SetupTest() {
	suite.translater.Original = make(map[string]string)
	suite.translater.Translation = make(map[string]string)
}

func (suite *YAMLTestSuite) TestTranslate() {
	suite.translater.Original = map[string]string{
		"greeting": "Hello world!",
	}
	err := suite.translater.Translate(language.English, language.French)
	suite.NoError(err)
	suite.Equal("Bonjour le monde!", suite.translater.Translation["greeting"])

	suite.translater.Client = nil
	err = suite.translater.Translate(language.English, language.French)
	suite.Error(err)
}

func (suite *YAMLTestSuite) TestParseFile() {
	err := suite.translater.ParseFile("../testfiles/demo.yml")
	suite.NoError(err)
	suite.Equal("Hello world!", suite.translater.Original["greeting"])

	err = suite.translater.ParseFile("./not-exists.yml")
	suite.Error(err)

	err = suite.translater.ParseFile("../testfiles/invalid.yml")
	suite.Error(err)
}

func (suite *YAMLTestSuite) TestSaveResult() {
	suite.translater.Original = map[string]string{
		"a": "a",
	}

	file := "/tmp/translate-result.yml"
	err := suite.translater.SaveResult(file)
	suite.NoError(err)
	_, err = os.Stat(file)
	suite.NoError(err)
	os.Remove(file)

	file = "/tmp/readonly.yml"
	ioutil.WriteFile(file, nil, 0444)
	err = suite.translater.SaveResult(file)
	suite.Error(err)
	os.Remove(file)
}
