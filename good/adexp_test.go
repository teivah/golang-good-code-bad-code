package good

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

// Test the parsing of a simple ADEXP message
func TestParseAdexpMessage(t *testing.T) {
	bytes, _ := ioutil.ReadFile("../resources/tests/adexp.txt")

	m, _ := ParseAdexpMessage(bytes)

	// Test upper level
	assert.Equal(t, true, m.IsUpperLevel())

	// Simple
	assert.Equal(t, "IFPL", m.Title)
	assert.Equal(t, "CYYZ", m.Adep)
	assert.Equal(t, "EASTERN :CREEK'()+,./", m.Alternate)
	assert.Equal(t, "AFIL", m.Ades)
	assert.Equal(t, "ACA878", m.Arcid)
	assert.Equal(t, "A333", m.ArcType)
	assert.Equal(t, "SDE3FGHIJ3J5LM1ORVWXY", m.Ceqpt)

	// Repeating
	assert.Equal(t, 13, len(m.Eetfir))
	assert.Equal(t, 2, len(m.Speed))

	// Complex
	assert.Equal(t, 2, len(m.Estdata))
	assert.Equal(t, 3, len(m.Geo))
	assert.Equal(t, 5, len(m.RoutePoints))

	// Route points
	assert.Equal(t, "CYYZ", m.RoutePoints[0].Ptid)
	assert.Equal(t, 0, m.RoutePoints[0].FlightLevel)
	assert.Equal(t, "170301220429", m.RoutePoints[0].Eto)
	assert.Equal(t, "JOOPY", m.RoutePoints[1].Ptid)
	assert.Equal(t, 390, m.RoutePoints[1].FlightLevel)
	assert.Equal(t, "170302002327", m.RoutePoints[1].Eto)
	assert.Equal(t, "GEO01", m.RoutePoints[2].Ptid)
	assert.Equal(t, 390, m.RoutePoints[2].FlightLevel)
	assert.Equal(t, "170302003347", m.RoutePoints[2].Eto)
	assert.Equal(t, "BLM", m.RoutePoints[3].Ptid)
	assert.Equal(t, 171, m.RoutePoints[3].FlightLevel)
	assert.Equal(t, "170302051642", m.RoutePoints[3].Eto)
	assert.Equal(t, "LSZH", m.RoutePoints[4].Ptid)
	assert.Equal(t, 14, m.RoutePoints[4].FlightLevel)
	assert.Equal(t, "170302052710", m.RoutePoints[4].Eto)
}

// Test ...
func TestParseSimpleToken(t *testing.T) {
	// Do something
}

// Test ...
func TestParseComplexToken(t *testing.T) {
	// Do something
}

// Performance test of an ADEXP message parsing
func BenchmarkParseAdexpMessage(t *testing.B) {
	log.SetLevel(log.FatalLevel)
	bytes, _ := ioutil.ReadFile("../resources/tests/adexp.txt")

	for i := 0; i < t.N; i++ {
		ParseAdexpMessage(bytes)
	}
}
