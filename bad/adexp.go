package bad

import (
	"strings"
	"fmt"
	"regexp"
	"sync"
)

const (
	patternSubfield = "-.[^-]*"

	tokenTitle   = "TITLE"
	tokenAdep    = "ADEP"
	tokenAltnz   = "ALTNZ"
	tokenAdes    = "ADES"
	tokenArcid   = "ARCID"
	tokenArctyp  = "ARCTYP"
	tokenCeqpt   = "CEQPT"
	tokenMsgtxt  = "MSGTXT"
	tokenComment = "COMMENT"
	tokenEetfir  = "EETFIR"
	tokenSpeed   = "SPEED"
	tokenEstdata = "ESTDATA"
	tokenGeo     = "GEO"
	tokenRtepts  = "RTEPTS"

	subtokenPtid   = "PTID"
	subtokenEto    = "ETO"
	subtokenFl     = "FL"
	subtokenGeoid  = "GEOID"
	subtokenLattd  = "LATTD"
	subtokenLongtd = "LONGTD"
)

var (
	stringNewline     = "\n"
	stringNewlineDash = "\n-"
	stringBegin       = "-BEGIN"
	stringEnd         = "-END"
	stringEmpty       = " "
	stringComment     = "//"
	stringDash        = "-"

	regexpSubfield, _ = regexp.Compile(patternSubfield)

	factory = map[string]func(string, string) interface{}{
		tokenTitle:   parseSimpleToken,
		tokenAdep:    parseSimpleToken,
		tokenAltnz:   parseSimpleToken,
		tokenAdes:    parseSimpleToken,
		tokenArcid:   parseSimpleToken,
		tokenArctyp:  parseSimpleToken,
		tokenCeqpt:   parseSimpleToken,
		tokenMsgtxt:  parseSimpleToken,
		tokenComment: parseSimpleToken,

		tokenEetfir: parseSimpleToken,
		tokenSpeed:  parseSimpleToken,

		tokenEstdata: parseComplexToken,
		tokenGeo:     parseComplexToken,
		tokenRtepts:  parseComplexToken,
	}

	mutexEetfir  = &sync.RWMutex{}
	mutexSpeed   = &sync.RWMutex{}
	mutexEstdata = &sync.RWMutex{}
	mutexGeo     = &sync.RWMutex{}
	mutexRtepts  = &sync.RWMutex{}
)

func ParseAdexpMessage(string string) (Message, error) {
	preprocessed, _ := preprocess(string)
	message := process(preprocessed)

	return message, nil
}

func preprocess(in string) (string, error) {
	if len(in) == 0 {
		return "", fmt.Errorf("Input is empty")
	}

	lines := strings.Split(in, stringNewline)
	var result string

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if startWith(line, stringEnd) {
		} else if startWith(line, stringBegin) {
			trimed := strings.Trim(line, stringEmpty)
			result = result + stringNewlineDash + trimed[len(stringBegin)+1:]
		} else if startWith(line, stringDash) {
			result = result + stringNewline + strings.Trim(line, stringEmpty)
		} else if startWith(line, stringEmpty) {
			result = result + stringEmpty + strings.Trim(line, stringEmpty)
		} else {
			result = result + stringEmpty + strings.Trim(line, stringEmpty)
		}
	}

	return result, nil
}

func process(in string) Message {
	lines := strings.Split(in, stringNewline)

	ch := make(chan string, len(lines))

	msg := Message{}

	for i := 0; i < len(lines); i++ {
		go mapLine(&msg, lines[i], ch)
	}

	for i := 0; i < len(lines); i++ {
		<-ch
	}

	msg.Type = AdexpType

	return msg
}

func mapLine(msg *Message, in string, ch chan string) {
	if !startWith(in, stringComment) {
		token, value := parseLine(in)
		if token != "" {
			f, contains := factory[string(token)]
			if !contains {
				ch <- "ok"
			} else {
				data := f(token, value)
				enrichMessage(msg, data)
				ch <- "ok"
			}
		} else {
			ch <- "ok"
			return
		}
	} else {
		ch <- "ok"
		return
	}
}

func enrichMessage(msg *Message, data interface{}) {
	if data == nil {
		return
	}

	switch data.(type) {
	case simpleToken:
		simpleToken := data.(simpleToken)
		value := simpleToken.value
		switch simpleToken.token {
		case tokenTitle:
			msg.Title = value
		case tokenAdep:
			msg.Adep = value
		case tokenAltnz:
			msg.Alternate = value
		case tokenAdes:
			msg.Ades = value
		case tokenArcid:
			msg.Arcid = value
		case tokenArctyp:
			msg.ArcType = value
		case tokenCeqpt:
			msg.Ceqpt = value
		case tokenMsgtxt:
			msg.MessageText = value
		case tokenComment:
			msg.Comment = value
		case tokenEetfir:
			mutexEetfir.Lock()
			msg.Eetfir = append(msg.Eetfir, value)
			mutexEetfir.Unlock()
		case tokenSpeed:
			mutexSpeed.Lock()
			msg.Speed = append(msg.Speed, value)
			mutexSpeed.Unlock()
		}
	case complexToken:
		complexToken := data.(complexToken)
		value := complexToken.value
		switch complexToken.token {
		case tokenEstdata:
			mutexEstdata.Lock()
			for _, v := range value {
				fl := extractFlightLevel(v[subtokenFl])
				msg.Estdata = append(msg.Estdata, estdata{v[subtokenPtid], v[subtokenEto], fl})
			}
			mutexEstdata.Unlock()
		case tokenGeo:
			mutexGeo.Lock()
			for _, v := range value {
				msg.Geo = append(msg.Geo, geo{v[subtokenGeoid], v[subtokenLattd], v[subtokenLongtd]})
			}
			mutexGeo.Unlock()
		case tokenRtepts:
			mutexRtepts.Lock()
			for _, v := range value {
				fl := extractFlightLevel(v[subtokenFl])
				msg.RoutePoints = append(msg.RoutePoints, rtepts{v[subtokenPtid], fl, v[subtokenEto]})
			}
			mutexRtepts.Unlock()
		}
	}
}

func parseSimpleToken(token, value string) interface{} {
	return simpleToken{string(token), string(value)}
}

func parseComplexToken(token, value string) interface{} {
	if value == "" {
		return complexToken{string(token), nil}
	}

	v := parseComplexLines(value, make(map[string]string), nil)

	return complexToken{string(token), v}
}

func parseComplexLines(in string, currentMap map[string]string, out []map[string]string) []map[string]string {
	match := regexpSubfield.Find([]byte(in))

	if match == nil {
		out = append(out, currentMap)
		return out
	}

	sub := string(match)

	h, l := parseLine(sub)

	_, contains := currentMap[string(h)]

	if contains {
		out = append(out, currentMap)
		currentMap = make(map[string]string)
	}

	currentMap[string(h)] = string(strings.Trim(l, stringEmpty))

	return parseComplexLines(in[len(sub):], currentMap, out)
}
