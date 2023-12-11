package utils

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

var currencyList = []string{"AFN", "ALL", "DZD", "AOA", "ARS", "AMD", "AWG", "AUD", "AZN", "BSD", "BHD", "BBD", "BDT", "BYN", "BZD", "BMD", "BTN", "XBT", "BOB", "BAM", "BWP", "BRL", "BND", "BGN", "BIF", "XPF", "KHR", "CAD", "CVE", "KYD", "FCFA", "CLP", "CLF", "CNY", "CNY", "COP", "CF", "CHF", "CDF", "CRC", "HRK", "CUC", "CZK", "DKK", "DJF", "DOP", "XCD", "EGP", "ETB", "FJD", "GMD", "GBP", "GEL", "GHS", "GTQ", "GNF", "GYD", "HTG", "HNL", "HKD", "HUF", "ISK", "INR", "IDR", "IRR", "IQD", "ILS", "JMD", "JPY", "JOD", "KMF", "KZT", "KES", "KWD", "KGS", "LAK", "LBP", "LSL", "LRD", "LYD", "MOP", "MKD", "MGA", "MWK", "MYR", "MVR", "MRO", "MUR", "MXN", "MDL", "MAD", "MZN", "MMK", "NAD", "NPR", "ANG", "NZD", "NIO", "NGN", "NOK", "OMR", "PKR", "PAB", "PGK", "PYG", "PHP", "PLN", "QAR", "RON", "RUB", "RWF", "SVC", "SAR", "RSD", "SCR", "SLL", "SGD", "SBD", "SOS", "ZAR", "KRW", "VES", "LKR", "SDG", "SRD", "SZL", "SEK", "CHF", "TJS", "TZS", "THB", "TOP", "TTD", "TND", "TRY", "TMT", "UGX", "UAH", "AED", "USD", "UYU", "UZS", "VND", "XOF", "YER", "ZMW", "ETH", "EUR", "LTC", "TWD", "PEN"}

type APIResponseStruct struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
	Error  interface{}            `json:"error"`
}

type CredentialStruct struct {
	ApiEndPoint string
	ApiKey      string
}

func GatewayResponseMapper(statusCode int, data map[string]interface{}, err interface{}) (events.APIGatewayProxyResponse, error) {
	body := APIResponseStruct{
		Status: "success",
		Data:   data,
		Error:  make(map[string]interface{}),
	}
	if err != nil {
		body = APIResponseStruct{
			Status: "failed",
			Data:   make(map[string]interface{}),
			Error:  err,
		}
	}
	jsonBody, _ := json.Marshal(body)
	log.Println("Response-->", string(jsonBody))

	response := events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"content-type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
		MultiValueHeaders: nil,
		Body:              string(jsonBody),
		IsBase64Encoded:   false,
	}
	return response, nil
}

func LoadAdapterCredentials(credential *CredentialStruct) {
	credential.ApiKey = os.Getenv("KEY")
	credential.ApiEndPoint = os.Getenv("API_ENDPOINT")
}

func IsCurrencyCodeValid(currencyCode string) bool {
	if !(len(strings.TrimSpace(currencyCode)) > 0) {
		return false
	}
	match, _ := regexp.MatchString(`^[A-Za-z]{3}$`, currencyCode)

	if !match {
		return false
	}
	if !slices.Contains(currencyList, strings.ToUpper(currencyCode)) {
		return false
	}
	return true
}
