package blockchain

const feePercentage = 0.01

func calculateFee(c Coin) Coin {
	fee := c.Float64() * feePercentage

	return ToCoin(fee)
}
