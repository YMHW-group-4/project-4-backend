package blockchain

const feePercentage = 0.01

// note: this is not (yet) implemented

func calculateFee(c Coin) Coin {
	fee := c.Float64() * feePercentage

	return ToCoin(fee)
}
