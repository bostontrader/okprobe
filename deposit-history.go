package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	//"reflect"

	//"crypto/hmac"
	//"crypto/sha256"
	//"encoding/base64"
	//"encoding/json"
	//"fmt"
	"github.com/bostontrader/okcommon"
	"io/ioutil"

	//"io/ioutil"
	"net/http"
	"time"
)

func ProbeDepositHistory(urlBase string, keyFile string) {

	endpoint := "/api/account/v3/deposit/history"
	url := urlBase + endpoint
	client := GetClient(urlBase)

	TestitStd(client, url)

	// In order to proceed we need to get real credentials.  Read them from a file.
	var obj utils.Credentials

	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Print(err)
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}

	// 7.
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30004(), 400) // OK-ACCESS-PASSPHRASE header is required

	// 8.
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30015(), 400) // Invalid OK_ACCESS_PASSPHRASE

	// 9.
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30013(), 401) // Invalid Sign

	// Requests after this point require a valid signature.

	// 11. No request params.  Final output.
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + endpoint
	encoded, _ := utils.HmacSha256Base64Signer(prehash, obj.SecretKey)
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)

	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
		"Vary":                      "",
	}
	body := Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	depositHistories := make([]utils.DepositHistory, 0)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&depositHistories)
	if err != nil {
		panic(err)
	}
	fmt.Println(&depositHistories)
	fmt.Println(reflect.TypeOf(depositHistories))
}
