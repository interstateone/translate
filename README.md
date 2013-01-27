# Bing Translate

A little wrapper for the Bing Translate API

## Usage

    import (
        "log"
        "translate"
    )

    myConfig = Config{
    	GrantType:    "client_credentials",
    	ScopeUrl:     "http://api.microsofttranslator.com",
    	ClientId:     "YourAppId",
    	ClientSecret: "YourClientSecret",
    	AuthUrl:      "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13/",
    }

    token := translate.GetToken(config)
    youreWhatTheFrenchCall := token.Translate("Les Incompétents", "fr", "en")
    log.Println(youreWhatTheFrenchCall)

Some simple tests are included and can be run with a `test_config.json` that looks like:

    {
      "ClientId": "YourAppID",
      "ClientSecret": "YourClientSecret"
    }


## License

Released under an MIT license, see the LICENSE file for more details.
