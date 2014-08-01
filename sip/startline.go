package sip

// Imports from the go standard library
import (
	"errors"
	"strings"
)

const (
	SIP_REQUEST          = "REQUEST"
	SIP_RESPONSE         = "RESPONSE"
	SIP_METHOD_INVITE    = "INVITE"
	SIP_METHOD_ACK       = "ACK"
	SIP_METHOD_OPTIONS   = "OPTIONS"
	SIP_METHOD_BYE       = "BYE"
	SIP_METHOD_CANCEL    = "CANCEL"
	SIP_METHOD_REGISTER  = "REGISTER"
	SIP_METHOD_INFO      = "INFO"
	SIP_METHOD_PRACK     = "PRACK"
	SIP_METHOD_SUBSCRIBE = "SUBSCRIBE"
	SIP_METHOD_NOTIFY    = "NOTIFY"
	SIP_METHOD_UPDATE    = "UPDATE"
	SIP_METHOD_MESSAGE   = "MESSAGE"
	SIP_METHOD_REFER     = "REFER"
	SIP_METHOD_PUBLISH   = "PUBLISH"
)

type StartLine struct {
	Body     string
	Type     string
	Method   string
	URI      string
	Resp     string
	RespText string
	Proto    string
	Version  string
}

func parseStartLine(startLine string) (*StartLine, error) {
	if len(startLine) < 3 {
		return nil, errors.New("parseStartLine err: length of start line is too short.")
	}
	s := &StartLine{}
	err := s.Parse(startLine)

	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *StartLine) Parse(sl string) error {
	s.Body = sl
	if sl[0:3] == "SIP" {
		s.Type = SIP_RESPONSE
		// Parse as response...
		parts := strings.SplitN(s.Body, " ", 3)
		if len(parts) != 3 {
			return errors.New("StartLine Parse Response err: err getting parts from LWS.")
		}
		charPos := strings.IndexRune(parts[0], '/')
		if charPos == -1 {
			return errors.New("StartLine Parse Response err: err getting proto char.")
		}
		s.Proto = parts[0][0:charPos]
		if len(parts[0])-1 < charPos+1 {
			return errors.New("StartLine Parse Response err: proto char appears to be at end of proto.")
		}
		s.Version = parts[0][charPos+1:]
		s.Resp = parts[1]
		s.RespText = parts[2]
		return nil
	}
	// Else, parse as request
	s.Type = SIP_REQUEST
	parts := strings.SplitN(s.Body, " ", 3)
	if len(parts) != 3 {
		return errors.New("StartLine Parse Request err: request line didn't split on LWS correctly.")
	}
	s.Method = parts[0]
	s.URI = parts[1]

	charPos := strings.IndexRune(parts[2], '/')
	if charPos == -1 {
		return errors.New("StartLine Parse Request err: could not get \"/\" pos in parts[2].")
	}
	if len(parts[2])-1 < charPos+1 {
		return errors.New("StartLine Parse Request err: \"/\" char appears to be at end of line.")
	}
	s.Proto = parts[2][0:charPos]
	s.Version = parts[2][charPos+1:]
	return nil

}
