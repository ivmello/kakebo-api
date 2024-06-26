package utils

import (
	"fmt"
)

func FormatMoney(value int, thousand, decimal string) string {
	var result string
	var isNegative bool
	if value < 0 {
		value = value * -1
		isNegative = true
	}
	result = fmt.Sprintf("%s%02d%s", decimal, value%100, result)
	value /= 100
	for value >= 1000 {
		result = fmt.Sprintf("%s%03d%s", thousand, value%1000, result)
		value /= 1000
	}
	if isNegative {
		return fmt.Sprintf("-%d%s", value, result)
	}
	return fmt.Sprintf("%d%s", value, result)
}
