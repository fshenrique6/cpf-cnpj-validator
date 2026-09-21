package validator

import "errors"

var ErrInvalidDocument = errors.New("documento inválido: não é um CPF nem um CNPJ válido")

func OnlyDigits(s string) string {
	return onlyDigits(s)
}

func DetectType(number string) (string, error) {
	number = onlyDigits(number)

	switch len(number) {
	case 11:
		if IsValidCPF(number) {
			return "cpf", nil
		}
	case 14:
		if IsValidCNPJ(number) {
			return "cnpj", nil
		}
	}

	return "", ErrInvalidDocument
}