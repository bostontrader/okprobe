package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"io/ioutil"
	"net/http"
	"time"
)

func ProbeAccountDepositAddress(urlBase string, keyFile string, makeErrors bool, query string) {

	endpoint := "/api/account/v3/deposit/address"
	url := urlBase + endpoint
	client := GetClient(urlBase)
	credentials := getCredentials(keyFile)

	if makeErrors {
		TestitStd(client, url, credentials, utils.ExpectedResponseHeaders)
	}

	// Requests after this point require a valid signature.

	// Extraneous parameters appear to be ignored, the full list is returned, and two extra headers appear: Vary and Strict-Transport-Security.  Status = 200.  For example: ?catfood and ?catfood=yum.
	// But this is of minimal importance so don't bother trying to test for this.  We certainly don't care to mimic this behavior in the catbox.

	if makeErrors {
		// Don't request any currency
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash := timestamp + "GET" + endpoint
		encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30023("currency cannot be blank"), 400)

		// Request an invalid currency
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		invalidParam := "catfood"
		params := "?currency=" + invalidParam
		prehash = timestamp + "GET" + endpoint + params
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
		req, _ = http.NewRequest("GET", url+params, nil)
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30031(invalidParam), 400)
	}

	// Try to submit a valid request by feeding a query string.  Even if we don't specify errors, we might still make them here.
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + endpoint + query
	encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
	req, err := http.NewRequest("GET", url+query, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Look for all of the expected headers in the received headers.
	compareHeaders(resp.Header, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), req, "deposit-address")

	if resp.StatusCode == 200 {

		// Parse this json just to prove that we can.
		depositAddress := make([]utils.DepositAddress, 0)
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		err = dec.Decode(&depositAddress)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(body))
		return
	} else if resp.StatusCode == 400 {
		fmt.Println(string(body))
		return
	} else {
		fmt.Println("StatusCode error:expected= 200 || 400, received=", resp.StatusCode)
		fmt.Println(string(body))
		return
	}

}
