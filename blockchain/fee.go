package blockchain

const feePercentage = 0.01

// TODO implement this

func calculateFee(c Coin) Coin {
	fee := c.Float64() * feePercentage

	return ToCoin(fee)
}
