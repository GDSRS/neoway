package utils

import (
	"fmt"
	"strings"
)

func GetInputLine(line []string) string {
	line[3] = handleDateField(line[3])
	//format money	
	line[4] = handleMoneyField(line[4])
	line[5] = handleMoneyField(line[5])

	return fmt.Sprintf("\n('%s',%t,%t,%s,%s,%s,'%s','%s')",
	line[0], line[1] == "1", line[2] == "1", line[3],
	line[4], line[5], line[6], line[7])
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