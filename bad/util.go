package bad

import (
	"strings"
	"strconv"
)

func startWith(in, test string) bool {
	if len(test) > len(in) {
		return false
	}

	i := 0
	for range test {
		if test[i] != in[i] {
			return false
		}
		i++
	}

	return true
}

func parseLine(in string) (string, string) {
	if len(in) == 0 {
		return "", ""
	}

	i := strings.Index(in, stringEmpty)

	if i == -1 {
		return in[1:], ""
	}

	return in[1:i], in[i+1:]
}

func extractFlightLevel(in string) int {
	fl, _ := strconv.Atoi(in[1:])
	return fl
}
