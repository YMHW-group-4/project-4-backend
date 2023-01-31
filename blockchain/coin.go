package blockchain

import "github.com/shopspring/decimal"

// Coin represents the currency within the blockchain.
type Coin struct {
	decimal decimal.Decimal
}

// ToCoin converts a float64 to Coin.
func ToCoin(f float64) Coin {
	return Coin{decimal.NewFromFloatWithExponent(f, -2)}
}

// Add adds a float64 to the Coin.
func (c Coin) Add(f float64) Coin {
	return Coin{c.decimal.Add(decimal.NewFromFloatWithExponent(f, -2))}
}

// Sub subtracts a float64 from the Coin.
func (c Coin) Sub(f float64) Coin {
	return Coin{c.decimal.Sub(decimal.NewFromFloatWithExponent(f, -2))}
}

// Float64 returns the Coin as a float64.
func (c Coin) Float64() float64 {
	f, _ := c.decimal.Float64()

	return f
}

// Equal checks if two coins are equal.
func (c Coin) Equal(coin Coin) bool {
	return c.decimal.Equal(coin.decimal)
}

// String returns a formatted Coin value.
func (c Coin) String() string {
	return c.decimal.String()
}
