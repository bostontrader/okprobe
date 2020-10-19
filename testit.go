package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"time"
)

// This is a standard sequence of errors related to credentials to submit to endpoints that are invoked via GET.
func TestitCredentialsErrors(httpClient *http.Client, url string, credentials utils.Credentials, expectedResponseHeaders map[string]string) error {

	// 1.
	req, _ := http.NewRequest("GET", url, nil)
	_, err := TestitAPI4xx(httpClient, req, 401, utils.ExpectedResponseHeaders, utils.Err30001())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 1: ", err)
		return err
	}

	// 2.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	_, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30002())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 2: ", err)
		return err
	}

	// 3.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	_, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30003())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 3: ", err)
		return err
	}

	// 4.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "invalid")
	_, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30005())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 4: ", err)
		return err
	}

	time.Sleep(1 * time.Second) // avoid rate limit

	// 5.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "2020-01-01T01:01:01.000Z") // expired
	_, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30008())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 5: ", err)
		return err
	}

	// 6. Set a good time stamp.  The system time is probably close enough to the server to work.  Maybe try to probe how far off the time can be.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	_, err = TestitAPI4xx(httpClient, req, 401, utils.ExpectedResponseHeaders, utils.Err30006())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 6: ", err)
		return err
	}

	// 7.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	_, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30004())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 7: ", err)
		return err
	}

	// 8.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
	_, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30015())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 8: ", err)
		return err
	}

	// 9.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	_, err = TestitAPI4xx(httpClient, req, 401, utils.ExpectedResponseHeaders, utils.Err30013())
	if err != nil {
		fmt.Println("TestitCredentials.  Error invoking API 9: ", err)
		return err
	}
	return nil
}

// This is a standard sequence of errors related to credentials to submit to endpoints that are invoked via POST.  There are enough differences between this function and TestitStd to justify the existence of this function.
func TestitStdPOST(client *http.Client, url string, credentials utils.Credentials, expectedResponseHeaders map[string]string) {

	// 1.
	req, _ := http.NewRequest("POST", url, nil)
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30001(), 401) // OK-ACCESS-KEY header is required

	// 2.
	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30002(), 400) // OK-ACCESS-SIGN header is required

	// 3.
	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30003(), 400) // OK-ACCESS-TIMESTAMP header is required

	// 4.
	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "invalid")
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30005(), 400) // Invalid OK-ACCESS-TIMESTAMP

	time.Sleep(1 * time.Second) // limit 6/sec

	// 5.
	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "2020-01-01T01:01:01.000Z")      // expired
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30008(), 400) // Request timestamp expired

	// 6. Set a good time stamp.  The system time is probably close enough to the server to work.  Maybe try to probe how far off the time can be.
	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30006(), 401) // Invalid OK-ACCESS-KEY

	// 7.
	req, _ = http.NewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30007(), 400) // Invalid Content_Type, please use the application/json format

	// Send body as nil, []byte(""), []byte(``) gives us a 500 error.

	// 8.
	body := []byte(`{}`)
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30004(), 400) // OK-ACCESS-PASSPHRASE header is required

	// 9.
	req, _ = http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
	Testit4xx(client, req, expectedResponseHeaders, utils.Err30015(), 400) // Invalid OK_ACCESS_PASSPHRASE

	// This test may or may not pass depending upon whether or not we are using the correct type of credentials.  Read only is no good.  Read and trade is probably ok.  How about read and withdraw?
	// 10.
	//req, _ = http.NewRequest("POST", url, bytes.NewBuffer(body))
	//req.Header.Add("Content-Type", "application/json")
	//req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	//req.Header.Add("OK-ACCESS-SIGN", "wrong")
	//req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	//req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	//Testit4xx(client, req, expectedResponseHeaders, utils.Err30013(), 401) // Invalid Sign

}

/*
Given an http client and a request, we want to make the API call and examine the response.  We want to ensure that
the we get the expected status, headers, and error messages, if applicable.  Read and close the response body, return
said body as a string, and return any error message if applicable.

This general process is confounded because status 2xx, 4xx, and 5xx are mostly the same but want to deal with error
messages using different types.  So...

TestitAPICore provides the common functionality...
TestitAPI2xx uses the core and does not expect any error.
TestitAPI4xx uses the core and expects a utils.OKError
TestitAPI5xx uses the core and expects a OK500Error
*/
func TestitAPICore(
	client *http.Client,
	req *http.Request,
	expectedStatusCode int,
	expectedResponseHeaders map[string]string,
) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error in TestitAPICore ", err)
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println("Error in TestitAPICore ", err)
		return []byte{}, err
	}

	if resp.StatusCode != expectedStatusCode {
		fmt.Println("Status code error: expected= ", expectedStatusCode, "received= ", resp.StatusCode, " body= ", string(body))
		return []byte{}, errors.New("Status code error.")
	}

	// Look for all of the expected headers in received headers.
	compareHeaders(resp.Header, expectedResponseHeaders, req)

	return body, nil
}

func TestitAPICoreNew(
	client *http.Client,
	req *http.Request,
	expectedStatusCode int,
) []byte {
	methodName := "okprobe:testit.go:TestitAPICoreNew"
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s: client.Do error %v\n", methodName, err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("%s: ioutil.ReadAll error %v\n", methodName, err)
		os.Exit(1)
	}

	resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		fmt.Printf("%s: Status code error: expected=%d, received=%d\nbody=%s\n", methodName, expectedStatusCode, resp.StatusCode, " body= ", string(body))
		os.Exit(1)
	}

	return body
}

func TestitAPI4xx(
	client *http.Client,
	req *http.Request,
	expectedStatusCode int,
	expectedResponseHeaders map[string]string,
	expectedErrorMessage utils.OKError,
) (string, error) {

	body, err := TestitAPICore(client, req, expectedStatusCode, expectedResponseHeaders)
	if err != nil {
		fmt.Println("Error in TestitAPI4xx")
		return "", err
	}

	var obj utils.OKError
	err = json.Unmarshal(body, &obj)
	if err != nil {
		fmt.Println("Error in TestitAPI4xx json.Unmarshal error: body=", string(body))
		return "", err
	}

	if !reflect.DeepEqual(obj, expectedErrorMessage) {
		fmt.Println("Error in TestitAPI4xx Error message compare error: expected=", expectedErrorMessage, ", received=", obj)
		return "", err
	}

	return string(body), nil

}

func TestitAPI4xxNew(
	client *http.Client,
	req *http.Request,
	expectedStatusCode int,
	expectedErrorMessage utils.OKError,
) {
	methodName := "okprobe:testit.go:TestitAPI4xxNew"
	body := TestitAPICoreNew(client, req, expectedStatusCode)

	var errorMessage utils.OKError
	err := json.Unmarshal(body, &errorMessage)
	if err != nil {
		fmt.Printf("%s: json.Unmarshal error: body=%s\n", methodName, string(body))
		os.Exit(1)
	}

	if !reflect.DeepEqual(errorMessage, expectedErrorMessage) {
		fmt.Printf("%s: Error message compare error: expected=%v, received=%v\n", methodName, expectedErrorMessage, errorMessage)
		os.Exit(1)
	}

	return
}

func TestitAPI5xx(
	client *http.Client,
	req *http.Request,
	expectedResponseHeaders map[string]string,
	expectedErrorMessage utils.OK500Error,
) (string, error) {

	body, err := TestitAPICore(client, req, 500, expectedResponseHeaders)
	if err != nil {
		fmt.Println("Error in TestitAPI5xx")
		return "", err
	}

	var obj utils.OK500Error
	err = json.Unmarshal(body, &obj)
	if err != nil {
		fmt.Println("Error in TestitAPI5xx json.Unmarshal error: body=", string(body))
		return "", err
	}

	if !reflect.DeepEqual(obj, expectedErrorMessage) {
		fmt.Println("Error in TestitAPI5xx Error message compare error: expected=", expectedErrorMessage, ", received=", obj)
		return "", err
	}

	return string(body), nil

}

func TestitAPI2xx(
	client *http.Client,
	req *http.Request,
	expectedResponseHeaders map[string]string,
) (string, error) {

	body, err := TestitAPICore(client, req, 200, expectedResponseHeaders)
	if err != nil {
		fmt.Println("Error in TestitAPI2xx")
		return "", err
	}

	return string(body), nil

}

/*
Deprecated. Use TestitAPI2xx. Given a client, and a request that we expect will produce a 200 response, make the API call and examine the response.  Ensure that
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
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("error:\nexpected= ", 200, "\nreceived=", resp.StatusCode)
		fmt.Println("Status code error: expected= ", 200, "received= ", resp.StatusCode, " body= ", string(body))
	}

	// Look for all of the expected headers in received headers.
	compareHeadersOld(resp.Header, expectedResponseHeaders, req, "tag200")

	// Read the body into a []byte and then create a return a new io.Reader using this []byte.  This enables us to close resp.Body, which we must do, and return an io.Reader which the caller needs in order to Decode JSON.
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(body)
}

/*
Deprecated. Use TestitAPI4xx. Given a client and a request, that we expect will produce a 4xx error, make the API call and examine the response.  Ensure that
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
	fmt.Println(string(body))
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
	compareHeadersOld(resp.Header, expectedResponseHeaders, req, "tag4xx")

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
			//fmt.Println(req.URL, key, " expected, but not present.")
		}
	}

	// Look for all of the received headers in expected headers.
	for key, _ := range respHeaders {
		_, ok := expectedResponseHeaders[key]
		if ok {
			// received and expected, cool
		} else {
			//fmt.Println(req.URL, key, " present, but not expected.")
		}
	}
}

// Did we receive all the headers we expected? Did we expect all that we received?
func compareHeadersOld(respHeaders map[string][]string, expectedResponseHeaders map[string]string, req *http.Request, tag string) {
	// Look for all of the expected headers in received headers.
	for key, _ := range expectedResponseHeaders {
		_, ok := respHeaders[key]
		if ok {
			// expected and present, cool
		} else {
			//fmt.Println(req.URL, tag, key, " expected, but not present.")
		}
	}

	// Look for all of the received headers in expected headers.
	for key, _ := range respHeaders {
		_, ok := expectedResponseHeaders[key]
		if ok {
			// received and expected, cool
		} else {
			//fmt.Println(req.URL, tag, key, " present, but not expected.")
		}
	}
}
