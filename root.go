package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	// For global flags
	baseURL                        string
	credentialsFile                string
	forReal                        bool
	makeErrorsCredentials          bool
	makeErrorsParams               bool
	makeErrorsWrongCredentialsType bool
	postBody                       string
	queryString                    string

	// For accountDepositHistoryByCurCmd flag
	currency string

	rootCmd = &cobra.Command{
		Use:   "okprobe",
		Short: "Work the OKEx API via the command line.",
		Long:  `Invoke the OKEx API, optionally using a variety of incorrect inputs in order to observe its behavior.`,
	}

	accountCurrenciesCmd = &cobra.Command{
		Use:   "accountCurrencies",
		Short: "Invoke GET /api/account/v3/currencies",
		Long:  `Invoke GET /api/account/v3/currencies`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountCurrencies(baseURL, credentialsFile, makeErrorsCredentials)
		},
	}

	accountDepositAddressCmd = &cobra.Command{
		Use:   "accountDepositAddress",
		Short: "Invoke GET /api/account/v3/deposit/address",
		Long:  `Invoke GET /api/account/v3/deposit/address`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountDepositAddress(baseURL, credentialsFile, makeErrorsCredentials, makeErrorsParams, queryString)
		},
	}

	accountDepositHistoryCmd = &cobra.Command{
		Use:   "accountDepositHistory",
		Short: "Invoke GET /api/account/v3/deposit/history",
		Long:  `Invoke GET /api/account/v3/deposit/history`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountDepositHistory(baseURL, credentialsFile, makeErrorsCredentials)
		},
	}

	accountDepositHistoryByCurCmd = &cobra.Command{
		Use:   "accountDepositHistoryByCur",
		Short: "Invoke GET /api/account/v3/deposit/history/<currency>",
		Long:  `Invoke GET /api/account/v3/deposit/history/<currency>`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountDepositHistoryByCur(baseURL, credentialsFile, makeErrorsCredentials, makeErrorsParams)
		},
	}

	accountLedgerCmd = &cobra.Command{
		Use:   "accountLedger",
		Short: "Invoke GET /api/account/v3/ledger",
		Long:  `Invoke GET /api/account/v3/ledger`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountLedger(baseURL, credentialsFile, makeErrorsCredentials, makeErrorsParams, queryString)
		},
	}

	accountTransferCmd = &cobra.Command{
		Use:   "accountTransfer",
		Short: "Invoke POST /api/account/v3/transfer",
		Long:  `Invoke POST /api/account/v3/transfer`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountTransfer(baseURL, credentialsFile, makeErrorsCredentials, makeErrorsParams, makeErrorsWrongCredentialsType, forReal, postBody)
		},
	}

	accountWalletCmd = &cobra.Command{
		Use:   "accountWallet",
		Short: "Invoke GET /api/account/v3/wallet",
		Long:  `Invoke GET /api/account/v3/wallet`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountWallet(baseURL, credentialsFile, makeErrorsCredentials)
		},
	}

	accountWithdrawalCmd = &cobra.Command{
		Use:   "accountWithdrawal",
		Short: "Invoke GET /api/account/v3/withdrawal",
		Long:  `Invoke GET /api/account/v3/withdrawal`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountWithdrawal(baseURL, credentialsFile, makeErrorsCredentials, makeErrorsParams, makeErrorsWrongCredentialsType, forReal, postBody)
		},
	}

	accountWithdrawalFeeCmd = &cobra.Command{
		Use:   "accountWithdrawalFee",
		Short: "Invoke GET /api/account/v3/withdrawal/fee",
		Long:  `Invoke GET /api/account/v3/withdrawal/fee`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeAccountWithdrawalFee(baseURL, credentialsFile, makeErrorsCredentials, makeErrorsParams, queryString)
		},
	}

	spotAccountsCmd = &cobra.Command{
		Use:   "spotAccounts",
		Short: "Invoke GET /api/spot/v3/accounts",
		Long:  `Invoke GET /api/spot/v3/accounts`,
		Run: func(cmd *cobra.Command, args []string) {
			ProbeSpotAccounts(baseURL, credentialsFile, makeErrorsCredentials)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of OKProbe",
		Long:  `All software has versions. This is OKProbe's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("OKProbe v0.1 -- HEAD")
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	//cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&baseURL, "baseURL", "https://www.okex.com", "The base url of the OKEx server to use.")
	rootCmd.PersistentFlags().StringVar(&credentialsFile, "credentialsFile", "/path/to/mycredentials.json", "A file name for a file that contains the OKEx credentials to use.")
	rootCmd.PersistentFlags().BoolVar(&forReal, "forReal", false, "When invoking this API send a real request that will actually produce a real response. (default false)")
	rootCmd.PersistentFlags().BoolVar(&makeErrorsCredentials, "makeErrorsCredentials", false, "When invoking this API send a variety errors related to the credentials and test that we receive the expected responses. (default false)")
	rootCmd.PersistentFlags().BoolVar(&makeErrorsParams, "makeErrorsParams", false, "When invoking this API send a variety errors related to the parameters of this endpoint and test that we receive the expected responses. (default false)")
	rootCmd.PersistentFlags().BoolVar(&makeErrorsWrongCredentialsType, "makeErrorsWrongCredentialsType", false, "When invoking this API, the user is knowingly sending the credentials of the wrong type.  Test that the call fails as expected. (default false)")
	rootCmd.PersistentFlags().StringVar(&postBody, "postBody", "", "A string to send to a POST endpoint for use as the request body.  For example: `{\"from\":\"6\", \"to\":\"1\", \"amount\":\"0.1\", \"currency\":\"bsv\"}`  (default empty string)")
	rootCmd.PersistentFlags().StringVar(&queryString, "queryString", "", "A complete query string to send to a GET endpoint.  For example: '?param1=A&param2=B'  (default empty string)")

	accountDepositHistoryByCurCmd.PersistentFlags().StringVar(&currency, "currency", "", "A currency of interest.'  (default empty string)")

	rootCmd.AddCommand(accountCurrenciesCmd)
	rootCmd.AddCommand(accountDepositAddressCmd)
	rootCmd.AddCommand(accountDepositHistoryCmd)
	rootCmd.AddCommand(accountDepositHistoryByCurCmd)
	rootCmd.AddCommand(accountLedgerCmd)
	rootCmd.AddCommand(accountTransferCmd)
	rootCmd.AddCommand(accountWalletCmd)
	rootCmd.AddCommand(accountWithdrawalCmd)
	rootCmd.AddCommand(accountWithdrawalFeeCmd)
	rootCmd.AddCommand(spotAccountsCmd)
	rootCmd.AddCommand(versionCmd)

}
