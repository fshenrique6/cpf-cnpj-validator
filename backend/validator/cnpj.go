package validator

func IsValidCNPJ(cnpj string) bool {
	cnpj = onlyDigits(cnpj)

	if len(cnpj) != 14 || allDigitsEqual(cnpj) {
		return false
	}

	digits := make([]int, 14)
	for i, c := range cnpj {
		digits[i] = int(c - '0')
	}

	firstCheckDigit := calculateCNPJCheckDigit(digits[:12])
	if firstCheckDigit != digits[12] {
		return false
	}

	secondCheckDigit := calculateCNPJCheckDigit(digits[:13])
	return secondCheckDigit == digits[13]
}

func calculateCNPJCheckDigit(digits []int) int {
	weights := make([]int, len(digits))

	if len(digits) == 12 {
		weights = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	} else {
		weights = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	}

	sum := 0
	for i, d := range digits {
		sum += d * weights[i]
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}