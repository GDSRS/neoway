package utils

import (
	"fmt"
	"unicode"
	"strings"
	"strconv"
	"errors"
	"log/slog"
)

func GetInputLine(line []string) (string, error) {
	// validate cpf

	if ok := line[0] == "NULL" || validateCpf(line[0]); !ok {
		return "", errors.New(fmt.Sprintf("Invalid CPF %s", line[0]))
	}
	//validate cnpj
	if ok := line[6] == "NULL" || validateCnpj(line[6]); !ok {
		return "", errors.New(fmt.Sprintf("Invalid CNPJ %s", line[6]))
	}
	if ok := line[7] == "NULL" || validateCnpj(line[7]); !ok {
		return "", errors.New(fmt.Sprintf("Invalid CNPJ %s", line[7]))
	}
	// format date field
	line[3] = handleDateField(line[3])
	// format money
	line[4] = handleMoneyField(line[4])
	line[5] = handleMoneyField(line[5])

	return fmt.Sprintf("\n('%s',%t,%t,%s,%s,%s,'%s','%s')",
		line[0], line[1] == "1", line[2] == "1", line[3],
		line[4], line[5], line[6], line[7]), nil
}

func validateCpf(cpf string) bool {
	cpf = getOnlyNumbers(cpf)
	if len(cpf) != 11 {
		return false
	}
	// validate first digit
	firstDigitHelper := []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	firstDigitSum, err := multiplyCPFCNPJWithHelpers(cpf, firstDigitHelper)
	if err != nil {
		return false
	}
	
	firstDigitResult := (firstDigitSum * 10) % 11
	slog.Debug(fmt.Sprintf("Total sum for first part is %d and first digit result is %d", firstDigitSum, firstDigitResult))
	if !(firstDigitResult == int(cpf[9] - '0') || (firstDigitResult == 10 && cpf[9] == '0')) {
		return false
	}

	//validate second digit
	secondDigitHelper := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
	secondDigitSum, err := multiplyCPFCNPJWithHelpers(cpf, secondDigitHelper)
	if err != nil {
		return false
	}
	secondDigitResult := (secondDigitSum * 10) % 11
	slog.Debug(fmt.Sprintf("Total sum for second part is %d and second digit result is %d, %d", secondDigitSum, secondDigitResult, int(cpf[10] - '0')))

	if !(secondDigitResult == int(cpf[10] - '0') || (secondDigitResult == 10 && cpf[10] == '0')) {
		return false
	}

	return true
}

func validateCnpj(cnpj string) bool {
	cnpj = getOnlyNumbers(cnpj)

	if len(cnpj) != 14 {
		return false
	}

	// validate first digit
	firstDigitHelper := []int{5,4,3,2,9,8,7,6,5,4,3,2}
	firstDigitSum, err := multiplyCPFCNPJWithHelpers(cnpj, firstDigitHelper)
	slog.Debug(fmt.Sprintf("Total sum for first part is %d", firstDigitSum))
	if err != nil {
		return false
	}

	firstDigitResult := firstDigitSum % 11
	if firstDigitResult < 2 {
		firstDigitResult = 0
	} else {
		firstDigitResult = 11 - firstDigitResult
	}

	if firstDigitResult != int(cnpj[12] - '0') {
		return false
	}

	// validate second digit
	secondDigitHelper := []int{6,5,4,3,2,9,8,7,6,5,4,3,2}
	secondDigitSum, err := multiplyCPFCNPJWithHelpers(cnpj, secondDigitHelper)
	slog.Debug(fmt.Sprintf("Total sum for second part is %d", secondDigitSum))
	if err != nil {
		return false
	}

	secondDigitResult := secondDigitSum % 11
	if secondDigitResult < 2 {
		secondDigitResult = 0
	} else {
		secondDigitResult = 11 - secondDigitResult
	}

	if secondDigitResult != int(cnpj[13] - '0') {
		return false
	}

	return true
}

func handleMoneyField(base string) string {
	if base == "NULL" {
		return base
	}

	var builder strings.Builder
	builder.Grow(len(base) + 2)

	builder.WriteRune('\'')
	for _, char := range base {
		if char == '.' {
			continue
		} else if char == ',' {
			builder.WriteRune('.')
		} else {
			builder.WriteRune(char)
		}
	}
	builder.WriteRune('\'')
	return builder.String()
}

func handleDateField(base string) string {
	if base != "NULL" {
		return fmt.Sprintf("'%s'", base)
	}
	return base
}

func getOnlyNumbers(base string) string {
	// return only numbers from a string
	var builder strings.Builder
	builder.Grow(len(base))
	for _, char := range(base) {
		if unicode.IsDigit(char) {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func multiplyCPFCNPJWithHelpers(base string, helper []int) (int, error) {
	totalSum := 0
	for i :=0; i < len(helper); i++ {
		num, err := strconv.Atoi(string(base[i]))
		if err != nil {
			slog.Error(fmt.Sprintf("Failed converting number %c: %v", base[i], err))
			return -1, err
		}

		totalSum += (num*helper[i])
	}
	return totalSum, nil
}
