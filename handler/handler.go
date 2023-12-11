package handler

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/opensaucerer/goaxios"
	"lambda-currency-converter-golang/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type APIResponseStruct struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
	Error  interface{}            `json:"error"`
}

type AdapterResponse struct {
	Quotes    map[string]float64 `json:"quotes"`
	Success   bool               `json:"success"`
	Source    string             `json:"source"`
	Timestamp float64            `json:"timestamp"`
}

type CacheMapItemStruct struct {
	Key       string
	Rate      float64
	CachedAt  time.Time
	ExpiresAt time.Time
}

var RateCache = map[string]CacheMapItemStruct{}
var CacheExpiry = 1 //in days
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fromCurrency := request.QueryStringParameters["fromCurrency"]
	toCurrency := request.QueryStringParameters["toCurrency"]
	amount, parseErr := strconv.ParseFloat(request.QueryStringParameters["amount"], 64)

	log.Println("Query::Params-->, Amount:: ", amount, "fromCurrency:: ", fromCurrency, "fromCurrency:: ", toCurrency)

	if parseErr != nil || amount <= 0 {
		return utils.GatewayResponseMapper(400, nil, "Invalid amount to perform FX")
	}

	if !utils.IsCurrencyCodeValid(toCurrency) {
		return utils.GatewayResponseMapper(400, nil, "Invalid 'toCurrency'")
	}

	if !utils.IsCurrencyCodeValid(fromCurrency) {
		return utils.GatewayResponseMapper(400, nil, "Invalid 'fromCurrency'")
	}

	credential := utils.CredentialStruct{}
	utils.LoadAdapterCredentials(&credential)

	currencyKey := strings.ToUpper(fromCurrency) + strings.ToUpper(toCurrency)
	var rate float64
	if RateCache[currencyKey].Rate != 0 && RateCache[currencyKey].ExpiresAt.Compare(time.Now()) == 1 {
		rate = RateCache[currencyKey].Rate
	} else {
		url := fmt.Sprintf("%s?access_key=%s&currencies=%s&source=%s",
			credential.ApiEndPoint, credential.ApiKey, strings.ToUpper(toCurrency), strings.ToUpper(fromCurrency))

		axiosCall := goaxios.GoAxios{
			Url:            url,
			Method:         "GET",
			ResponseStruct: &AdapterResponse{},
		}

		response := axiosCall.RunRest()

		if response.Error != nil {
			log.Println("Adapter Error-->", response.Error)
			return utils.GatewayResponseMapper(500, nil, "Unknown error occurred, please try again later")
		}

		parsedData, _ := response.Body.(*AdapterResponse)

		rate = parsedData.Quotes[currencyKey]
		if rate == 0 {
			return utils.GatewayResponseMapper(500, nil, "Failed to fetch fx rate, please try again later")
		}
		RateCache[currencyKey] = CacheMapItemStruct{
			Key:       currencyKey,
			Rate:      rate,
			CachedAt:  time.Now().UTC(),
			ExpiresAt: time.Now().UTC().AddDate(0, 0, CacheExpiry),
		}
	}

	return utils.GatewayResponseMapper(200, map[string]interface{}{
		"actualAmount":    amount,
		"convertedAmount": rate * amount,
		"rate":            rate,
		"fromCurrency":    fromCurrency,
		"toCurrency":      toCurrency,
	}, nil)
}
