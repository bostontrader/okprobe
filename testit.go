package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

/*
Every API call has a standard sequence of errors submitted.
*/
func TestitStd(client *http.Client, url string, credentials utils.Credentials) {

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

	// 7.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30004(), 400) // OK-ACCESS-PASSPHRASE header is required

	// 8.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30015(), 400) // Invalid OK_ACCESS_PASSPHRASE

	// 9.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30013(), 401) // Invalid Sign

}

/*
Given a client, and a request that we expect will produce a 200 response, make the API call and examine the response.  Ensure that
the we get the expected headers and response.
*/
func Testit200(
	client *http.Client,
	req *http.Request,
	expectedResponseHeaders map[string]string,
) io.Reader {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("error:\nexpected= ", 200, "\nreceived=", resp.StatusCode)
	}

	// Look for all of the expected headers in received headers.
	compareHeaders(resp.Header, expectedResponseHeaders, req)

	// Read the body into a []byte and then create a return a new io.Reader using this []byte.  This enables us to close resp.Body, which we must do, and return an io.Reader which the caller needs in order to Decode JSON.
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(body)
}

/*
Given a client and a request, that we expect will produce a 4xx error, make the API call and examine the response.  Ensure that
the we get the expected headers, error message, and status code.
*/
func Testit4xx(
	client *http.Client,
	req *http.Request,
	expectedResponseHeaders map[string]string,
	expectedErrorMessage utils.OKError,
	expectedStatusCode int,
) {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != expectedStatusCode {
		fmt.Println("error:\nexpected= ", expectedStatusCode, "\nreceived=", resp.StatusCode)
	}

	var obj utils.OKError
	err = json.Unmarshal(body, &obj)
	if err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(obj, expectedErrorMessage) {
		fmt.Println("error:\nexpected=", expectedErrorMessage, "\nreceived=", obj)
	}

	// Look for all of the expected headers in received headers.
	compareHeaders(resp.Header, expectedResponseHeaders, req)

	return
}

// Did we receive all the headers we expected? Did we expect all that we received?
func compareHeaders(respHeaders map[string][]string, expectedResponseHeaders map[string]string, req *http.Request) {
	// Look for all of the expected headers in received headers.
	for key, _ := range expectedResponseHeaders {
		_, ok := respHeaders[key]
		if ok {
			// expected and present, cool
		} else {
			fmt.Println("KeyE2:", key, " expected, but not present.")
		}
	}

	// Look for all of the received headers in expected headers.
	for key, _ := range respHeaders {
		_, ok := expectedResponseHeaders[key]
		if ok {
			// received and expected, cool
		} else {
			fmt.Println("KeyR2:", key, " present, but not expected.")
		}
	}
}
