package main

import (
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

/* Given a urlBase, endPoint, and queryString, as well as credentials, built a suitable GET request.
If present, the queryString should include the initial ? mark. In the event of any errors, os.Exit(1)
*/
func buildGETRequest(credentials utils.Credentials, endPoint, queryString, urlBase string) *http.Request {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + endPoint + queryString
	encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
	req, err := http.NewRequest("GET", urlBase+endPoint+queryString, nil)
	if err != nil {
		fmt.Printf("okprobe:main.go:buildGETRequest: Error building NewRequest %v\n", err)
		os.Exit(1)
	}
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

	return req
}

/* Given a urlBase, endPoint, and postBody, as well as credentials, built a suitable JSON POST request.
In the event of any errors, os.Exit(1)
*/
func buildPOSTRequest(credentials utils.Credentials, endPoint, postBody, urlBase string) *http.Request {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "POST" + endPoint + postBody
	encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
	req, err := http.NewRequest("POST", urlBase+endPoint+queryString, nil)
	if err != nil {
		fmt.Printf("okprobe:main.go:buildPOSTRequest: Error building NewRequest %v\n", err)
		os.Exit(1)
	}
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	req.Header.Add("Content-Type", "application/json")

	return req
}

/* Given a file name, attempt to read and parse a credentials file.  os.exit(1) on any error. */
func getCredentials(keyFile string) utils.Credentials {
	methodName := "okprobe:main.go:getCredentials"
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Printf("%s: Error reading file %s, error=%v\n", methodName, keyFile, err)
		os.Exit(1)
	}

	var obj utils.Credentials
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Printf("%s: Error parsing credentials file, error=%v\ndata=%s", keyFile, err, string(data))
		os.Exit(1)
	}

	return obj
}

func main() {
	err := Execute()
	if err != nil {
		fmt.Printf("okprobe:main.go:main: Execute error %v\n", err)
		os.Exit(1)
	}
}
