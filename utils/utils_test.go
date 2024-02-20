package utils

import (
	"testing"
)

func TestValidateCpfTrue(t *testing.T) {
	cpf := "348.508.300-30"
	result := validateCpf(cpf)
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("validateCpf should have return %t for cpf %s but it returned %t", expectedResult, cpf, result)
	}
}

func TestValidateCpfFalse(t *testing.T) {
	cpf := "348.508.300-39"
	result := validateCpf(cpf)
	expectedResult := true
	if result != expectedResult {
		t.Fatalf("validateCpf should have return %t for cpf %s but it returned %t", expectedResult, cpf, result)
	}
}