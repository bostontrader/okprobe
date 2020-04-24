package utils

type OKError struct {
	ErrorMessage string `json:"error_message"`
	Code         int    `json:"code"`
	ErrorCode    string `json:"error_code"`
	Message      string `json:"message"`
}

func Err30001() OKError {
	return OKError{
		ErrorMessage: "OK-ACCESS-KEY header is required",
		Code:         30001,
		ErrorCode:    "30001",
		Message:      "OK-ACCESS-KEY header is required",
	}
}

func Err30002() OKError {
	return OKError{
		ErrorMessage: "OK-ACCESS-SIGN header is required",
		Code:         30002,
		ErrorCode:    "30002",
		Message:      "OK-ACCESS-SIGN header is required",
	}
}

func Err30003() OKError {
	return OKError{
		ErrorMessage: "OK-ACCESS-TIMESTAMP header is required",
		Code:         30003,
		ErrorCode:    "30003",
		Message:      "OK-ACCESS-TIMESTAMP header is required",
	}
}

func Err30004() OKError {
	return OKError{
		ErrorMessage: "OK-ACCESS-PASSPHRASE header is required",
		Code:         30004,
		ErrorCode:    "30004",
		Message:      "OK-ACCESS-PASSPHRASE header is required",
	}
}

func Err30005() OKError {
	return OKError{
		ErrorMessage: "Invalid OK-ACCESS-TIMESTAMP",
		Code:         30005,
		ErrorCode:    "30005",
		Message:      "Invalid OK-ACCESS-TIMESTAMP",
	}
}

func Err30006() OKError {
	return OKError{
		ErrorMessage: "Invalid OK-ACCESS-KEY",
		Code:         30006,
		ErrorCode:    "30006",
		Message:      "Invalid OK-ACCESS-KEY",
	}
}

func Err30008() OKError {
	return OKError{
		ErrorMessage: "Request timestamp expired",
		Code:         30008,
		ErrorCode:    "30008",
		Message:      "Request timestamp expired",
	}
}

func Err30013() OKError {
	return OKError{
		ErrorMessage: "Invalid Sign",
		Code:         30013,
		ErrorCode:    "30013",
		Message:      "Invalid Sign",
	}
}

func Err30015() OKError {
	return OKError{
		ErrorMessage: "Invalid OK_ACCESS_PASSPHRASE",
		Code:         30015,
		ErrorCode:    "30015",
		Message:      "Invalid OK_ACCESS_PASSPHRASE",
	}
}
