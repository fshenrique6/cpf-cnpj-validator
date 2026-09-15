package validator

import "testing"

func TestIsValidCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{"CPF válido sem formatação", "11144477735", true},
		{"CPF válido com formatação", "111.444.777-35", true},
		{"CPF com todos os dígitos iguais", "11111111111", false},
		{"CPF com tamanho incorreto", "123456789", false},
		{"CPF com dígito verificador incorreto", "11144477736", false},
		{"CPF vazio", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCPF(tt.cpf)
			if result != tt.expected {
				t.Errorf("IsValidCPF(%q) = %v; esperado %v", tt.cpf, result, tt.expected)
			}
		})
	}
}