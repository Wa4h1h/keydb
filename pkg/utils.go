package pkg

import (
	"fmt"
	"strings"
)

func ParseStringToSlice(str string) ([]string, error) {
	if len(str) < 3 {
		return nil, ErrMalformedSlice
	}

	elements := str[1 : len(str)-1]

	if !(str[0] == '[') || !(str[len(str)-1] == ']') {
		return nil, ErrMalformedSlice
	}

	return strings.Split(elements, ","), nil
}

func ParseSliceToString(slice []string) string {
	str := strings.Join(slice, ",")

	return fmt.Sprintf("[%s]", str)
}