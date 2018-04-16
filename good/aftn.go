/*
Package good is a library for parsing the ADEXP messages.
An intermediate format Message is built by the parser.
*/

package good

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// These parts should be configurable obviously but for the sake of this exercise we will let them hardcoded.
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

// AFTN message types
const (
	AdexpType = iota
	IcaoType  = iota
)

// Upper level constant
const upperLevel = 350

// Estimated data
type estdata struct {
	Ptid        string // Point id
	Eto         string // Estimated Time Over
	FlightLevel int    // Flight level
}

// Geo point
type geo struct {
	Geoid     string // Geo point id
	Latitude  string // Point latitude
	Longitude string // Point longitude
}

// Route points
type rtepts struct {
	Ptid        string // Point id
	FlightLevel int    // Flight level
	Eto         string // Estimated Time Over
}

// Message structure produced by the parser
type Message struct {
	Type        int       // Message type (ADEXP or ICAO)
	Title       string    // Message title
	Adep        string    // Departure airport
	Ades        string    // Destination airport
	Alternate   string    // Alternate aerodrome
	Arcid       string    // Aircraft identifier
	ArcType     string    // Aircraft type
	Ceqpt       string    // Equipment
	MessageText string    // Message text
	Comment     string    // Personal comments
	Eetfir      []string  // Flight information region
	Speed       []string  // Speed
	Estdata     []estdata // Estimated data
	Geo         []geo     // Geo points
	RoutePoints []rtepts  // Route points
}

// Structure returned for parsing simple tokens
type simpleToken struct {
	token string
	value string
}

// Structure returned for parsing complex tokens
type complexToken struct {
	token string
	value []map[string]string
}

// Checks whether a message concerns the upper level (FL350)
func (m *Message) IsUpperLevel() bool {
	for _, r := range m.RoutePoints {
		if r.FlightLevel > upperLevel {
			return true
		}
	}

	return false
}
