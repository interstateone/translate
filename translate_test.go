package translate

import (
	"encoding/json"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type TranslateSuite struct {
	Config Config
}

var _ = Suite(&TranslateSuite{})

func (s *TranslateSuite) SetUpSuite(c *C) {
	s.Config = Config{
		GrantType:    "client_credentials",
		ScopeUrl:     "http://api.microsofttranslator.com",
		ClientId:     "",
		ClientSecret: "",
		AuthUrl:      "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13/",
	}

	filename := "config.json"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		c.FailNow()
	}
	err = json.Unmarshal(body, &s.Config)
	if err != nil {
		c.FailNow()
	}
}

func (s *TranslateSuite) TestGetToken(c *C) {
	_, err := GetToken(&s.Config)
	c.Assert(err, IsNil)
}

func (s *TranslateSuite) TestTranslate(c *C) {
	token, err := GetToken(&s.Config)

	german, err := token.Translate("", "", "de")
	c.Assert(german, Equals, "")
	c.Assert(err.Error(), Equals, "\"text\" is a required parameter")

	german, err = token.Translate("Black cats", "", "")
	c.Assert(german, Equals, "")
	c.Assert(err.Error(), Equals, "\"to\" is a required parameter")

	french, err := token.Translate("Purple centipedes", "en", "fr")
	c.Logf("French: %s", french)
	c.Assert(err, IsNil)

	spanish, err := token.Translate("Orange iguanas", "", "es")
	c.Logf("Spanish: %s", spanish)
	c.Assert(err, IsNil)
}
