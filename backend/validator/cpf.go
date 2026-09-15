package validator

import (
	"regexp"
	"strconv"
)

var nonDigitRegex = regexp.MustCompile(`[^0-9]`)

func onlyDigits(s string) string {
	return nonDigitRegex.ReplaceAllString(s, "")
}

func allDigitsEqual(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func IsValidCPF(cpf string) bool {
	cpf = onlyDigits(cpf)

	if len(cpf) != 11 || allDigitsEqual(cpf) {
		return false
	}

	digits := make([]int, 11)
	for i, c := range cpf {
		digits[i], _ = strconv.Atoi(string(c))
	}

	firstCheckDigit := calculateCPFCheckDigit(digits[:9], 10)
	if firstCheckDigit != digits[9] {
		return false
	}

	secondCheckDigit := calculateCPFCheckDigit(digits[:10], 11)
	return secondCheckDigit == digits[10]
}

func calculateCPFCheckDigit(digits []int, initialWeight int) int {
	sum := 0
	weight := initialWeight

	for _, d := range digits {
		sum += d * weight
		weight--
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}