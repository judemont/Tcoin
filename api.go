package main

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
)

const BASE_API_URL string = "https://openapiv1.coinstats.app"

var API_KEYS = [1]string{"vohN0pzG3mjbLsVYXFjyL+soUlp5tnqvuJQo7X74sxY="}

func getCoinData(cointId string) (Crypto, error) {

	url := BASE_API_URL + "/coins/" + cointId

	response := apiCall(url)

	var responseMap map[string]interface{}

	json.Unmarshal(response, &responseMap)

	if statusCode, ok := responseMap["statusCode"]; ok && statusCode != http.StatusOK {
		return Crypto{}, errors.New(responseMap["message"].(string))
	}

	var crypto Crypto = Crypto{
		ID:                       responseMap["id"].(string),
		Name:                     responseMap["name"].(string),
		Symbol:                   responseMap["symbol"].(string),
		Price:                    responseMap["price"].(float64),
		LogoURL:                  responseMap["icon"].(string),
		PriceChangePercentageDay: responseMap["priceChange1d"].(float64),
		Website:                  responseMap["websiteUrl"].(string),
		MarketCap:                responseMap["marketCap"].(float64),
		MarketCapRank:            int(responseMap["rank"].(float64)),
		TotalSupply:              responseMap["totalSupply"].(float64),
		CirculatingSupply:        responseMap["availableSupply"].(float64),
		Volume:                   responseMap["volume"].(float64),
	}
	return crypto, nil

}

func getChart(coinId string, period string) ([][]float64, error) {
	url := BASE_API_URL + "/coins/" + coinId + "/charts?period=" + period

	response := apiCall(url)

	var responseMap map[string]interface{}
	json.Unmarshal(response, &responseMap)

	if statusCode, ok := responseMap["statusCode"]; ok && statusCode != http.StatusOK {
		return nil, errors.New(responseMap["message"].(string))
	}

	var chartData [][]float64
	if err := json.Unmarshal(response, &chartData); err != nil {
		return nil, err
	}
	return chartData, nil
}

func apiCall(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-KEY", API_KEYS[rand.Intn(len(API_KEYS))])

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	return body
}
