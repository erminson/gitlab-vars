package utils

import (
	"errors"
	"fmt"
	"strings"
)

func ProcessArgs(args []string) string {
	newArs := make([]string, 0)
	newArs = append(newArs, args...)

	jsonStr := strings.Join(newArs, "")
	jsonStr = strings.ReplaceAll(jsonStr, " ", "")
	return strings.TrimSpace(jsonStr)
}

func ProcessString(raw string) (string, error) {
	if !strings.HasPrefix(raw, "{") || !strings.HasSuffix(raw, "}") {
		return "", errors.New("no correct JSON")
	}
	raw = raw[1 : len(raw)-1]

	oldSlice := strings.Split(raw, ",")
	newSlice := make([]string, 0)
	for _, pairRaw := range oldSlice {
		pairSlice := strings.Split(pairRaw, ":")
		if len(pairSlice) != 2 {
			return "", errors.New("no correct JSON")
		}

		key := fmt.Sprintf("\"%s\"", pairSlice[0])
		var value interface{}
		if pairSlice[1] == "false" || pairSlice[1] == "true" {
			value = fmt.Sprintf("%v", pairSlice[1])
		} else {
			value = fmt.Sprintf("\"%v\"", pairSlice[1])
		}

		newSlice = append(newSlice, fmt.Sprintf("%s: %v", key, value))
	}

	return fmt.Sprintf("{%s}", strings.Join(newSlice, ", ")), nil
}
