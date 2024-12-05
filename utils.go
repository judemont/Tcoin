package main

import (
	"fmt"
	"strings"
)

func FormatPrice(symbol string, price float64) string {
	priceStr := fmt.Sprintf("%f", price)
	intPart, decPart := splitPrice(priceStr)
	intPart = addCommas(intPart)
	return fmt.Sprintf("%s.%s %s", intPart, decPart, symbol)
}

func splitPrice(price string) (string, string) {
	parts := strings.Split(price, ".")
	return parts[0], parts[1]
}

func addCommas(intPart string) string {
	n := len(intPart)
	if n <= 3 {
		return intPart
	}
	var result strings.Builder
	for i, digit := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune('\'')
		}
		result.WriteRune(digit)
	}
	return result.String()
}
