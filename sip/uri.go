package sip

import (
	"errors"
	//"strings"
)

const (
	SIP_SCHEME  = "sip"
	SIPS_SCHEME = "sips"
	TEL_SCHEME  = "tel"
)

type URI struct {
	Body     string
	Scheme   string
	User     string
	Password string
	Host     string
	Port     string
}

func parseURI(uri string) (*URI, error) {
	if len(uri) < 3 {
		return nil, errors.New("parseURI err: length of uri string is too short: " + uri)
	}
	newUri := &URI{}
	err := newUri.Parse(uri)

	if err != nil {
		return nil, err
	}
	return newUri, nil
}

func (u *URI) Parse(uriString) error {
	u.Body = uriString
	sLen := len(uriString)
	if uriString[0:4] == "sip:" {
		uriString = uriString[4:]
		u.Scheme = SIP_SCHEME
	}	
    else if u.Raw[0:4] == "tel:" {
		uriString = uriString[4:]
		u.Scheme = TEL_SCHEME
	}
	else if sLen > 5 && u.Raw[0:5] == "sips:" {
        uriString = uriString[5:]
		u.Scheme = SIPS_SCHEME
	}
    else {
        return errors.New("parseURI err: Bad SIP URI. Must start with 'sip:', 'tel:', 'sips:'.")
    }
	
}
