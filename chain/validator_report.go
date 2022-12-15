package chain

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type ValidatorReport struct {
	Active            bool
	DelegatedTokens   string
	PendingCommission string
	Price             float64
	CalculatedApr     float64
	EstimatedApr      float64
	// TODO: Rank,
}

func (c *Chain) GetValidatorReport(validatorAddress string) (ValidatorReport, error) {
	querier := c.getQuerier()

	stakingValidatorResp, err := querier.Staking_Validator(validatorAddress)
	if err != nil {
		return ValidatorReport{}, err
	}

	distValidatorCommissionResp, err := querier.Distribution_ValidatorCommission(validatorAddress)
	if err != nil {
		return ValidatorReport{}, err
	}
	pendingCommission, err := distValidatorCommissionResp.Commission.Commission[0].Amount.Float64()
	if err != nil {
		return ValidatorReport{}, err
	}

	cosmosDirData, err := getCosmosDirectoryDataForChain(c.Name)

	return ValidatorReport{
		Active:            stakingValidatorResp.Validator.Status == stakingtypes.Bonded,
		DelegatedTokens:   c.getTokens(float64(stakingValidatorResp.Validator.Tokens.Int64())),
		PendingCommission: c.getTokens(pendingCommission),
		Price:             c.Price,
		CalculatedApr:     cosmosDirData.Chain.Params.CalculatedApr,
		EstimatedApr:      cosmosDirData.Chain.Params.EstimatedApr,
	}, nil
}
