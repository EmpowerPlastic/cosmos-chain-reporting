package main

import (
	"github.com/EmpowerPlastic/cosmos-chain-reporting/chain"
	"github.com/EmpowerPlastic/cosmos-chain-reporting/config"
	"go.uber.org/zap"
)

func main() {
	rootLogger, _ := zap.NewProduction()
	defer rootLogger.Sync() // flushes buffer, if any

	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		panic(err)
	}

	for _, chainCfg := range cfg.Chains {
		c, err := chain.LoadChain(rootLogger, chainCfg)
		if err != nil {
			panic(err)
		}

		if err := c.SetToLatestHeight(rootLogger); err != nil {
			panic(err)
		}

		/*if chainCfg.Validator != "" {
			validatorReport, err := c.GetValidatorReport(chainCfg.Validator)
			if err != nil {
				panic(err)
			}

			rootLogger.Sugar().Infoln(validatorReport)
		}*/

		/*for _, wallet := range chainCfg.Wallets {
			txs, err := c.GetTransactions(wallet.Address)
			if err != nil {
				panic(err)
			}

			_ = txs
		}*/


	}
}
