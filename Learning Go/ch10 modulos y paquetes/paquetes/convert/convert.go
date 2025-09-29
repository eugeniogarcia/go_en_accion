// Package convert provides various utilities to
// make it easy to convert money from one currency to another.
package convert

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Money represents the combination of an amount of money
// and the currency the money is in.
//
// The value is stored using a [github.com/shopspring/decimal.Decimal]
type Money struct {
	Value    decimal.Decimal
	Currency string
}

// Convert converts the value of one currency to another.
//
// It has two parameters: a Money instance with the value to convert,
// and a string that represents the currency to convert to. Convert returns
// the converted currency and any errors encountered from unknown or unconvertible
// currencies.
//
// If an error is returned, the Money instance is set to the zero value.
//
// Supported currencies are:
//   - USD - US Dollar
//   - CAD - Canadian Dollar
//   - EUR - Euro
//   - INR - Indian Rupee
//
// More information on exchange rates can be found at [Investopedia].
//
// [Investopedia]: https://www.investopedia.com/terms/e/exchangerate.asp
func Convert(from Money, to string) (Money, error) {
	switch to {
	case "USD":
		return Money{
			Value:    from.Value, // 1 USD = 1 USD
			Currency: "USD",
		}, nil
	case "CAD":
		return Money{
			Value:    from.Value.Mul(decimal.NewFromFloat(1.25)), // 1 USD = 1.25 CAD
			Currency: "CAD",
		}, nil
	case "EUR":
		return Money{
			Value:    from.Value.Mul(decimal.NewFromFloat(0.90)), // 1 USD = 0.90 EUR
			Currency: "EUR",
		}, nil
	case "INR":
		return Money{
			Value:    from.Value.Mul(decimal.NewFromFloat(74.0)), // 1 USD = 74.0 INR
			Currency: "INR",
		}, nil
	default:
		return Money{}, fmt.Errorf("unsupported currency: %s", to)
	}
}
