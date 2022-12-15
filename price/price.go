package price

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type responseJson []priceJson

type priceJson struct {
	Price float64
}

func GetPrice(denom string) (float64, error) {
	endpoint := fmt.Sprintf("https://api-osmosis.imperator.co/tokens/v2/%s", denom)
	resp, err := http.Get(endpoint)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("GET %s returned code %d", endpoint, resp.StatusCode)
	}

	resBody, err := io.ReadAll(resp.Body)
	var respJson responseJson
	if err := json.Unmarshal(resBody, &respJson); err != nil {
		return 0, err
	}
	if len(respJson) != 1 {
		return 0, fmt.Errorf("1 price was not found, instead found %d for denom %s", len(respJson), denom)
	}

	if respJson[0].Price == 0 {
		return 0, fmt.Errorf("price not found for denom %s", denom)
	}

	return respJson[0].Price, nil
}
