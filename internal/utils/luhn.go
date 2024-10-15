package utils

import (
	"strconv"
	"unicode"
)

func IsLuhn(orderId string) bool {
	var sum int
	var alternate bool

	for i := len(orderId) - 1; i >= 0; i-- {
		r := rune(orderId[i])

		if !unicode.IsDigit(r) {
			return false
		}

		n, _ := strconv.Atoi(string(r))

		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}
