package main

import (
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

func ProbeCurrencies(urlBase string, keyFile string) {

	endpoint := "/api/account/v3/currencies"
	url := urlBase + endpoint

	c1 := GetClient(urlBase)
	client := &c1

	// 1.
	req, _ := http.NewRequest("GET", url, nil)
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30001(), 401) // OK-ACCESS-KEY header is required

	// 2.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30002(), 400) // OK-ACCESS-SIGN header is required

	// 3.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30003(), 400) // OK-ACCESS-TIMESTAMP header is required

	// 4.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "invalid")
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30005(), 400) // Invalid OK-ACCESS-TIMESTAMP

	time.Sleep(1 * time.Second) // limit 6/sec

	// 5.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "2020-01-01T01:01:01.000Z")            // expired
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30008(), 400) // Request timestamp expired

	// 6. Set a good time stamp.  The system time is probably close enough to the server to work.  Maybe try to probe how far off the time can be.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30006(), 401) // Invalid OK-ACCESS-KEY

	// In order to proceed we need to get real credentials.  Read them from a file.
	type APIKey struct {
		Key        string `json:"api_key"`
		SecretKey  string `json:"api_secret_key"`
		Passphrase string `json:"passphrase"`
	}
	var obj APIKey
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Print(err)
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}

	// 7.
	req, err = http.NewRequest("GET", url, nil)
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

	// The final correct request must have a valid signature, so build one now.
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
	currenciesEntries := make([]utils.CurrenciesEntry, 0)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&currenciesEntries)
	if err != nil {
		panic(err)
	}
	fmt.Println(&currenciesEntries)
	fmt.Println(reflect.TypeOf(currenciesEntries))
}
