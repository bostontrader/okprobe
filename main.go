package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	endPtr := flag.String("endpnt", "wallet", "The name of the endpoint to invoke")
	keyFilePtr := flag.String("keyfile", "/path/to/apikeys.json", "The name of a file that contains the API keys")
	makeErrorsPtr := flag.Bool("errors", false, "Invoke the API calls using intentional errors.")

	urlPtr := flag.String("url", "https://www.okex.com", "The URL of the API")

	if len(os.Args) < 2 {
		flag.Usage()
	}

	flag.Parse()
	fmt.Println("Invoking the Probe with the following configuration:")
	fmt.Println("endpnt:", *endPtr)
	fmt.Println("keyfile:", *keyFilePtr)
	fmt.Println("errors:", *makeErrorsPtr)
	fmt.Println("url:", *urlPtr)

	switch *endPtr {
	case "currencies":
		ProbeCurrencies(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	case "deposit-address":
		ProbeDepositAddress(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	case "deposit-history":
		ProbeDepositHistory(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	case "wallet":
		ProbeWallet(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	case "withdrawal-fee":
		ProbeWithdrawalFee(*urlPtr, *keyFilePtr, *makeErrorsPtr)

	default:
		fmt.Println("Unknown endpoint ", *endPtr)
	}
}
