package bad

const (
	AdexpType = 0
	IcaoType  = 1
)

const upperLevel = 350

type estdata struct {
	Ptid        string
	Eto         string
	FlightLevel int
}

type geo struct {
	Geoid     string
	Latitude  string
	Longitude string
}

type rtepts struct {
	Ptid        string
	FlightLevel int
	Eto         string
}

type Message struct {
	Type        int
	Title       string
	Adep        string
	Ades        string
	Alternate   string
	Arcid       string
	ArcType     string
	Ceqpt       string
	MessageText string
	Comment     string
	Eetfir      []string
	Speed       []string
	Estdata     []estdata
	Geo         []geo
	RoutePoints []rtepts
}

type simpleToken struct {
	token string
	value string
}

type complexToken struct {
	token string
	value []map[string]string
}

func IsUpperLevel(m Message) bool {
	for _, r := range m.RoutePoints {
		if r.FlightLevel > upperLevel {
			return true
		}
	}

	return false
}
