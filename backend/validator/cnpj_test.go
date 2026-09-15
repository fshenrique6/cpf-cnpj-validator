package validator

import "testing"

func TestIsValidCNPJ(t *testing.T) {
	tests := []struct {
		name     string
		cnpj     string
		expected bool
	}{
		{"CNPJ válido sem formatação", "11222333000181", true},
		{"CNPJ válido com formatação", "11.222.333/0001-81", true},
		{"CNPJ com todos os dígitos iguais", "11111111111111", false},
		{"CNPJ com tamanho incorreto", "1122233300018", false},
		{"CNPJ com dígito verificador incorreto", "11222333000180", false},
		{"CNPJ vazio", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCNPJ(tt.cnpj)
			if result != tt.expected {
				t.Errorf("IsValidCNPJ(%q) = %v; esperado %v", tt.cnpj, result, tt.expected)
			}
		})
	}
}