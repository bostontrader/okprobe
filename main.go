package main

import (
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"io/ioutil"
	"net/http"
	"time"
)

// Build and execute an http request.  Test for a 200 status and the expected response headers.  Return the response body as a string.
func standardGET(client http.Client, credentials utils.Credentials, endPoint string, queryString string, urlBase string) string {

	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
		"Vary":                      "",
	}

	req, err := standardGETReq(credentials, endPoint, queryString, urlBase)
	if err != nil {
		fmt.Println("Error building the request ", err)
		return ""
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error in submitting the http request: ", err)
		return ""
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Status code error: expected= ", 200, " received=", resp.StatusCode, " body =", string(body))
		return ""
	}

	// Look for all of the expected headers in received headers.
	compareHeaders(resp.Header, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), req)

	// Read the body into a []byte and return it.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body. ", err)
		return ""
	}

	return string(body)
}

// queryString should have the initial ? mark, if present
func standardGETReq(credentials utils.Credentials, endPoint string, queryString, urlBase string) (*http.Request, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + endPoint + queryString
	encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
	req, err := http.NewRequest("GET", urlBase+endPoint+queryString, nil)
	if err != nil {
		fmt.Println("Error building NewRequest ", err)
		return nil, err
	}
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

	return req, nil
}

func getCredentials(keyFile string) (utils.Credentials, error) {
	var obj utils.Credentials
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Error reading file ", keyFile, err)
		return utils.Credentials{}, err
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("Error parsing credentials file ", keyFile, err)
		return utils.Credentials{}, err
	}

	return obj, nil
}

func getCredentialsOld(keyFile string) utils.Credentials {
	var obj utils.Credentials
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func catMap(a, b map[string]string) map[string]string {
	var n = map[string]string{}
	for k, v := range a {
		n[k] = v
	}
	for k, v := range b {
		n[k] = v
	}

	return n
}

func main() {
	Execute()
}
