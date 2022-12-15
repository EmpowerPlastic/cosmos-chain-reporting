package chain

import (
	"context"
	"fmt"
	"github.com/EmpowerPlastic/cosmos-chain-reporting/config"
	"github.com/EmpowerPlastic/cosmos-chain-reporting/price"
	"github.com/pkg/errors"
	"github.com/strangelove-ventures/lens/client"
	lens "github.com/strangelove-ventures/lens/client"
	registry "github.com/strangelove-ventures/lens/client/chain_registry"
	"github.com/strangelove-ventures/lens/client/query"
	"go.uber.org/zap"
	"math"
	"math/big"
	"os"
	"time"
)

type Chain struct {
	Name      string
	Denom     string
	Exponent  uint
	RestAPI   string
	Height    int64
	BlockTime time.Time
	Price     float64

	chainClient *client.ChainClient
}

func LoadChain(log *zap.Logger, chainConfig config.ChainConfig) (Chain, error) {
	chainInfo, err := registry.DefaultChainRegistry(log).GetChain(context.Background(), chainConfig.RegistryName)
	if err != nil {
		return Chain{}, errors.Wrap(err, "failed to get chain info")
	}

	ccc, err := chainInfo.GetChainConfig(context.Background())
	if err != nil {
		return Chain{}, errors.Wrap(err, "failed to get lens client chain config")
	}

	chainClient, err := lens.NewChainClient(log, ccc, os.Getenv("HOME"), os.Stdin, os.Stdout)
	if err != nil {
		return Chain{}, errors.Wrapf(err, "failed to build new chain client for %s", chainConfig.RegistryName)
	}

	price, err := price.GetPrice(chainConfig.Denom)
	if err != nil {
		return Chain{}, errors.Wrapf(err, "failed to get chain price for %s", chainConfig.RegistryName)
	}

	return Chain{
		Name:        chainConfig.RegistryName,
		Denom:       chainConfig.Denom,
		Exponent:    chainConfig.Exponent,
		RestAPI:     chainConfig.RestAPI,
		Price:       price,
		chainClient: chainClient,
	}, nil
}

func (c *Chain) SetToLatestHeight(log *zap.Logger) error {
	querier := query.Query{
		Client:  c.chainClient,
		Options: query.DefaultOptions(),
	}

	status, err := querier.Status()
	if err != nil {
		return err
	}

	c.Height = status.SyncInfo.LatestBlockHeight
	c.BlockTime = status.SyncInfo.LatestBlockTime
	log.Sugar().Infof("block height for chain %s set to %d", c.Name, c.Height)

	return nil
}

func (c *Chain) getQuerier() query.Query {
	options := query.DefaultOptions()
	options.Height = c.Height
	querier := query.Query{
		Client:  c.chainClient,
		Options: options,
	}

	return querier
}

func (c *Chain) getTokens(utokens float64) string {
	divideBy := math.Pow(10, float64(c.Exponent))
	tokens := new(big.Float).Quo(big.NewFloat(utokens), big.NewFloat(divideBy))
	tokensF, _ := tokens.Float64()
	return fmt.Sprintf("%f %s", tokensF, c.Denom)
}
