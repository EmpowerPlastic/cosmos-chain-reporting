package chain

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
)

type Tx struct {
	Hash      string
	Timestamp string
	Messages  []Message
}

type Message struct {
	Type       string
	From       string
	To         string
	FromAmount float64
	ToAmount   float64
	MsgRawData map[string]interface{}
}

type txsResponseJson struct {
	TxResponses []struct {
		Txhash    string
		Code      int
		Timestamp string
		Tx        struct {
			Body struct {
				Messages []map[string]interface{}
			}
		}
		Logs Logs
	} `json:"tx_responses"`
}

type Logs []struct {
	Events []Event
}

type Event struct {
	Type       string
	Attributes []Attribute
}

type Attribute struct {
	Key   string
	Value string
}

func (c *Chain) GetTransactions(address string) ([]Tx, error) {
	events := []string{
		fmt.Sprintf("message.sender%%3D%%27%s%%27", address),
		fmt.Sprintf("transfer.recipient%%3D%%27%s%%27", address),
		fmt.Sprintf("fungible_token_packet.receiver%%3D%%27%s%%27", address),
	}

	mappedTxs := make(map[string]Tx)
	for _, event := range events {
		endpoint := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs?events=%s", c.RestAPI, event)
		resp, err := http.Get(endpoint)
		if err != nil {
			return nil, err
		}
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "not able to ready body from GET %s with code %d", endpoint, resp.StatusCode)
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("GET %s returned code %d and body %s", endpoint, resp.StatusCode, string(resBody))
		}

		var respJson txsResponseJson
		if err := json.Unmarshal(resBody, &respJson); err != nil {
			return nil, err
		}

		for _, txResp := range respJson.TxResponses {
			if _, ok := mappedTxs[txResp.Txhash]; ok {
				continue
			}

			var messages []Message
			for i, m := range txResp.Tx.Body.Messages {
				t := m["@type"].(string)
				switch t {
				case "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission":
					e, err := txResp.Logs.getEvent(i, "withdraw_commission")
					if err != nil {
						return nil, errors.Wrap(err, txResp.Txhash)
					}
					amountAttr, err := e.findAttributeValue("amount")
					if err != nil {
						return nil, errors.Wrap(err, txResp.Txhash)
					}

					amount, err := strconv.ParseFloat(amountAttr, 64)

					messages = append(messages, Message{
						Type:       t,
						From:       "",
						To:         address,
						FromAmount: amount,
						ToAmount:   amount,
						MsgRawData: m,
					})
					// TODO: Look for withdraw_commission event (will need index)
				case "":

				}
			}

			mappedTxs[txResp.Txhash] = Tx{
				Hash:      txResp.Txhash,
				Timestamp: txResp.Timestamp,
				Messages:  nil,
			}
		}
	}

	return nil, nil
}

// TODO: Add support for optional list of attributes that need to match
func (l Logs) getEvent(messageIndex int, eventType string) (Event, error) {
	for _, event := range l[messageIndex].Events {
		if event.Type == eventType {
			return event, nil
		}
	}

	return Event{}, fmt.Errorf("not able to find eventType %q", eventType)
}

func (e Event) findAttributeValue(key string) (string, error) {
	for _, attr := range e.Attributes {
		if attr.Key == key {
			return attr.Value, nil
		}
	}

	return "", fmt.Errorf("not able to find attribute %q", key)
}