package chain

type BalanceReport struct {

}

func (c *Chain) GetBalanceReport(address string) (BalanceReport, error) {
	querier := c.getQuerier()
	balances, err := querier.Bank_Balances(address)
	if err != nil {
		return BalanceReport{}, err
	}
}