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

type HTMLTestSuite struct {
	suite.Suite
	translater *HTMLTranslater
}

func TestHTMLTestSuite(t *testing.T) {
	suite.Run(t, new(HTMLTestSuite))
}

// SetupSuite - run by testify once at the very start of the testing suite, before any tests are run.
func (suite *HTMLTestSuite) SetupSuite() {
	client, err := translate.NewClient(
		context.TODO(),
		option.WithAPIKey(os.Getenv("GCLOUD_API_KEY")),
	)
	suite.NoError(err)
	suite.translater = NewHTMLTranslater(client)
}

// TearDownSuite - run by testify once, at the very end of the testing suite, after all tests have been run.
func (suite *HTMLTestSuite) TearDownSuite() {
	suite.translater.Client.Close()
}

// SetupTest - run before every test in the suite.
func (suite *HTMLTestSuite) SetupTest() {
	suite.translater.Original = ""
	suite.translater.Translation = ""
}

func (suite *HTMLTestSuite) TestTranslate() {
	suite.translater.Original = "<h1>Hello world!</h1>"
	err := suite.translater.Translate(language.English, language.French)
	suite.NoError(err)
	suite.Equal("<h1> Bonjour le monde! </h1>", suite.translater.Translation)

	suite.translater.Client = nil
	err = suite.translater.Translate(language.English, language.French)
	suite.Error(err)
}

func (suite *HTMLTestSuite) TestParseFile() {
	err := suite.translater.ParseFile("../testfiles/demo.html")
	suite.NoError(err)
	suite.True(len(suite.translater.Original) > 0)

	err = suite.translater.ParseFile("./not-exists.html")
	suite.Error(err)
}

func (suite *HTMLTestSuite) TestSaveResult() {
	suite.translater.Translation = "<h1>Hello world!</h1>"

	file := "/tmp/translate-result.html"
	err := suite.translater.SaveResult(file)
	suite.NoError(err)
	_, err = os.Stat(file)
	suite.NoError(err)
	os.Remove(file)

	file = "/tmp/readonly.html"
	ioutil.WriteFile(file, nil, 0444)
	err = suite.translater.SaveResult(file)
	suite.Error(err)
	os.Remove(file)
}
