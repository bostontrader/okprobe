package main

import (
	"encoding/json"
	//"github.com/spf13/cobra/cobra/cmd"

	//"encoding/json"
	//"flag"
	//utils "github.com/bostontrader/okcommon"
	//"github.com/spf13/cobra"
	//"io/ioutil"
	//"os"
	//"flag"
	utils "github.com/bostontrader/okcommon"
	"io/ioutil"
	//"os"
	//"./cmd"
)

func main() {
	Execute()
}

func getCredentials(keyFile string) utils.Credentials {
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
