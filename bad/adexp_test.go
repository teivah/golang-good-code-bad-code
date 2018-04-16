package bad

import (
	"io/ioutil"
	"testing"
)

func BenchmarkParseAdexpMessage(t *testing.B) {
	bytes, _ := ioutil.ReadFile("../resources/tests/adexp.txt")

	for i := 0; i < t.N; i++ {
		ParseAdexpMessage(string(bytes))
	}
}
