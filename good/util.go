package good

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// Check whether a byte slice starts with another byte slice
func startWith(in, test []byte) bool {
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

// Parse a line by returning the header (token name) and the value. Example: -COMMENT TEST must returns COMMENT and TEST (in byte slices)
func parseLine(in []byte) ([]byte, []byte) {
	if len(in) == 0 {
		return nil, nil
	}

	i := bytes.Index(in, bytesEmpty)

	if i == -1 {
		return in[1:], nil
	}

	return in[1:i], in[i+1:]
}

// As each flight level in represented as Fxxx (e.g. F350), this function simply parse the flight level to an int with an eventual error (if the int conversion fails for instance)
func extractFlightLevel(in string) (int, error) {
	fl, err := strconv.Atoi(in[1:])

	if err != nil {
		log.Errorf("flight level %v cannot be parsed", fl)
		return -1, fmt.Errorf("flight level %v cannot be parsed", fl)
	}

	return fl, err
}
