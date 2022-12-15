package chain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CosmosDirectoryData struct {
	Chain struct {
		Params struct {
			CalculatedApr float64 `json:"calculated_apr"`
			EstimatedApr  float64 `json:"estimated_apr"`
		}
	}
}

func getCosmosDirectoryDataForChain(chain string) (CosmosDirectoryData, error) {
	endpoint := fmt.Sprintf("https://chains.cosmos.directory/%s", chain)
	resp, err := http.Get(endpoint)
	if err != nil {
		return CosmosDirectoryData{}, err
	}
	if resp.StatusCode != 200 {
		return CosmosDirectoryData{}, fmt.Errorf("GET %s returned code %d", endpoint, resp.StatusCode)
	}

	resBody, err := io.ReadAll(resp.Body)
	var respJson CosmosDirectoryData
	if err := json.Unmarshal(resBody, &respJson); err != nil {
		return CosmosDirectoryData{}, err
	}

	return respJson, nil
}
