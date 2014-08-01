package sip

import (
	"strconv"
	"strings"
)

const (
	HDR_ACCEPT                        = "accept"          // RFC3261
	HDR_ACCEPT_CONTACT                = "accept-contact"  // RFC3841
	HDR_ACCEPT_CONTACT_CMP            = "a"               // RFC3841
	HDR_ACCEPT_ENCODING               = "accept-encoding" //
	HDR_ACCEPT_LANGUAGE               = "accept-language"
	HDR_ACCEPT_RESOURCE_PRIORITY      = "accept-resource-priority" // RFC4412
	HDR_ALERT_INFO                    = "alert-info"
	HDR_ALLOW                         = "allow"
	HDR_ALLOW_EVENTS                  = "allow-events"
	HDR_ALLOW_EVENTS_CMP              = "u"
	HDR_ANSWER_MODE                   = "answer-mode"
	HDR_AUTHENTICATION_INFO           = "authentication-info"
	HDR_AUTHORIZATION                 = "authorization"
	HDR_CALL_ID                       = "call-id"
	HDR_CALL_ID_CMP                   = "i"
	HDR_CALL_INFO                     = "call-info"
	HDR_CONTACT                       = "contact"
	HDR_CONTACT_CMP                   = "m"
	HDR_CONTENT_DISPOSITION           = "content-disposition"
	HDR_CONTENT_ENCODING              = "content-encoding"
	HDR_CONTENT_ENCODING_CMP          = "e"
	HDR_CONTENT_LANGUAGE              = "content-language"
	HDR_CONTENT_LENGTH                = "content-length"
	HDR_CONTENT_LENGTH_CMP            = "l"
	HDR_CONTENT_TYPE                  = "content-type"
	HDR_CONTENT_TYPE_CMP              = "c"
	HDR_CSEQ                          = "cseq"
	HDR_DATE                          = "date"
	HDR_ERROR_INFO                    = "error-info"
	HDR_EVENT                         = "event"
	HDR_EXPIRES                       = "expires"
	HDR_FLOW_TIMER                    = "flow-timer"
	HDR_FROM                          = "from"
	HDR_FROM_CMP                      = "f"
	HDR_HISTORY_INFO                  = "history-info"  // from RFC 4244
	HDR_IDENTITY                      = "identity"      // RFC 4474
	HDR_IDENTITY_CMP                  = "y"             // RFC 4474
	HDR_IDENTITY_INFO                 = "identity-info" // RFC 4474
	HDR_IDENTITY_INFO_CMP             = "n"             // RFC 4474
	HDR_IN_REPLY_TO                   = "in-reply-to"
	HDR_JOIN                          = "join" // RFC 3911
	HDR_MAX_FORWARDS                  = "max-forwards"
	HDR_MIME_VERSION                  = "mime-version"
	HDR_MIN_EXPIRES                   = "min-expires"
	HDR_MIN_SE                        = "min-se" // RFC4028
	HDR_ORGANIZATION                  = "organization"
	HDR_PATH                          = "path"               // RFC3327
	HDR_PERMISSION_MISSING            = "permission-missing" // RFC5360
	HDR_PRIORITY                      = "priority"
	HDR_PRIVACY                       = "privacy"
	HDR_PRIV_ANSWER_MODE              = "priv-answer-mode" // RFC 5373
	HDR_PROXY_AUTHENTICATE            = "proxy-authenticate"
	HDR_PROXY_AutHORIZATION           = "proxy-authorization"
	HDR_PROXY_REQUIRE                 = "proxy-require"
	HDR_RACK                          = "rack" // RFC 3262
	HDR_REASON                        = "reason"
	HDR_RECORD_ROUTE                  = "record-route"
	HDR_REFER_SUB                     = "refer-sub"                     // RFC4488
	HDR_REFER_TO                      = "refer-to"                      // RFC 3515, RFC 4508
	HDR_REFERRED_BY                   = "referred-by"                   // RFC3892
	HDR_REFERRED_BY_CMP               = "b"                             // RFC3892
	HDR_REJECT_CONTACT                = "reject-contact"                // RFC3841
	HDR_REJECT_CONTACT_CMP            = "j"                             // RFC3841
	HDR_REMOTE_PARTY_ID               = "remote-party-id"               // DRAFT
	HDR_REPLACES                      = "replaces"                      // RFC3891
	HDR_REPLY_TO                      = "reply-to"                      // RFC3261
	HDR_REQUEST_DISPOSITION           = "request-disposition"           // RFC3841
	HDR_REQUIRE                       = "require"                       // RFC3261
	HDR_RESOURCE_PRIORITY             = "resource-priority"             // RFC4412
	HDR_RETRY_AFTER                   = "retry-after"                   // RFC3261
	HDR_ROUTE                         = "route"                         // RFC3261
	HDR_RSEQ                          = "rseq"                          // RFC3262
	HDR_SECUTIRY_CLIENT               = "security-client"               // RFC3329
	HDR_SECURITY_SERVER               = "security-server"               // RFC3329
	HDR_SECURITY_VERIFY               = "security-verify"               // RFC3329
	HDR_SERVER                        = "server"                        // RFC3261
	HDR_SERVICE_ROUTE                 = "service-route"                 // RFC3608
	HDR_SESSION_EXPIRES               = "session-expires"               // RFC4028
	HDR_SESSION_EXPIRES_CMP           = "x"                             // RFC4028
	HDR_ETAG                          = "sip-etag"                      // RFC3903
	HDR_IF_MATCH                      = "sip-if-match"                  // RFC3903
	HDR_SUBJECT                       = "subject"                       // RFC3261
	HDR_SUBJECT_CMP                   = "s"                             // RFC3261
	HDR_SUBSCRIPTION_STATE            = "subscription-state"            // RFC3265
	HDR_SUPPORTED                     = "supported"                     // RFC3261
	HDR_SUPPORTED_CMP                 = "k"                             // RFC3261
	HDR_SUPPRESS_IF_MATCH             = "suppress-if-match"             // RFC5839
	HDR_TARGET_DIALOG                 = "target-dialog"                 // RFC4538
	HDR_TIMESTAMP                     = "timestamp"                     // RFC3261
	HDR_TO                            = "to"                            // RFC3261
	HDR_TO_CMP                        = "t"                             // RFC3261
	HDR_TRIGGER_CONSENT               = "trigger-consent"               // RFC5360
	HDR_UNSUPPORTED                   = "unsupported"                   // RFC3261
	HDR_USER_AGENT                    = "user-agent"                    // RFC3261
	HDR_VIA                           = "via"                           // RFC3261
	HDR_VIA_CMP                       = "v"                             // RFC3261
	HDR_WARNING                       = "warning"                       // RFC3261
	HDR_WWW_AUTHENTICATE              = "www-authenticate"              // RFC3261
	HDR_P_ACCESS_NETWORK_INFO         = "p-access-network-info"         // RFC3455
	HDR_P_ANSWER_STATE                = "p-answer-state"                // RFC3455
	HDR_P_ASSERTED_IDENTITY           = "p-asserted-identity"           // RFC3325
	HDR_P_ASSERTED_SERVICE            = "p-asserted-service"            // RFC3455
	HDR_P_ASSOCIATED_URI              = "p-associated-uri"              // RFC3455
	HDR_P_CALLED_PARTY_ID             = "p-called-party-id"             // RFC3455
	HDR_P_CHARGING_FUNCTION_ADDRESSES = "p-charging-function-addresses" // RFC3455
	HDR_P_CHARGING_VECTOR             = "p-charging-vector"             // RFC3455
	HDR_P_DCS_BILLING_INFO            = "p-dcs-billing-info"            // RFC5503
	HDR_P_DCS_LAES                    = "p-dcs-laes"                    // RFC5503
	HDR_P_DCS_OSPS                    = "p-dcs-osps"                    // RFC5503
	HDR_P_DCS_REDIRECT                = "p-dcs-redirect"                // RFC5503
	HDR_P_DCS_TRACE_PARTY_ID          = "p-dcs-trace-party-id"          // RFC5503
	HDR_P_EARLY_MEDIA                 = "p-early-media"                 // RFC5009
	HDR_P_MEDIA_AUTHORIZATION         = "p-media-authorization"         // RFC3313
	HDR_P_PREFERRED_IDENTITY          = "p-preferred-identity"          // RFC3325
	HDR_P_PREFERRED_SERVICE           = "p-preferred-service"           // RFC6050
	HDR_P_PROFILE_KEY                 = "p-profile-key"                 // RFC5002
	HDR_P_USER_DATABASE               = "p-user-database"               // RFC4457
	HDR_P_VISITED_NETWORK_ID          = "p-visited-network-id"          // RFC3455
)

type SipMsg struct {
	State              string
	Error              error
	Msg                string
	Body               string
	StartLine          *StartLine
	Accept             string
	AlertInfo          string
	Allow              []string
	AllowEvents        []string
	Authorization      string
	ContentDisposition string
	ContentLength      string
	ContentLengthInt   int
	ContentType        string
	From               string
	MaxForwards        int64
	Organization       string
	To                 string
	Contact            string
	ContactVal         string
	CallId             string
	Cseq               string
	Rack               string
	Reason             string
	Rseq               string
	RseqInt            int
	RecordRoute        []string
	Route              []string
	Via                []string
	Require            []string
	Supported          []string
	Privacy            string
	ProxyAuthenticate  string
	ProxyRequire       []string
	RemotePartyIdVal   string
	RemotePartyId      string
	PAssertedIdVal     string
	PAssertedId        string
	Unsupported        []string
	UserAgent          string
	Server             string
	Subject            string
	WWWAuthenticate    string
	Warning            string
	OtherHeaders       []string
}

func NewSipMsg(packet []byte) (*SipMsg, error) {
	newSipMsg := &SipMsg{}
	newSipMsg.Parse(string(packet))
	return newSipMsg, nil
}

// Parse takes a SIP message and parses it into the SipMsg struct
func (s *SipMsg) Parse(m string) error {
	// Store the original message into the message property
	s.Msg = m

	// Split the message by newlines for individual headers
	headers := strings.Split(m, "\n")

	for i := range headers {
		if i == 0 {
			s.StartLine, _ = parseStartLine(headers[i])
		} else {
			if headers[i] == "\n" || headers[i] == "" {
				return nil
			}
			headerKeyValue := strings.SplitN(headers[i], ":", 2)
			headerKey := strings.TrimSpace(strings.ToLower(headerKeyValue[0]))
			headerValue := strings.TrimSpace(headerKeyValue[1])

			switch headerKey {
			case HDR_ACCEPT:
				s.Accept = headerValue
			case HDR_ALLOW:
				s.Allow = append(s.Allow, headerValue)
			case HDR_ALLOW_EVENTS:
				s.AllowEvents = append(s.AllowEvents, headerValue)
			case HDR_ALLOW_EVENTS_CMP:
				s.AllowEvents = append(s.AllowEvents, headerValue)
			case HDR_AUTHORIZATION:
				s.Authorization = headerValue
			case HDR_CALL_ID:
				s.CallId = headerValue
			case HDR_CALL_ID_CMP:
				s.CallId = headerValue
			case HDR_CONTACT:
				s.ContactVal = headerValue
			case HDR_CONTACT_CMP:
				s.ContactVal = headerValue
			case HDR_CONTENT_DISPOSITION:
				s.ContentDisposition = headerValue
			case HDR_CONTENT_LENGTH:
				s.ContentLength = headerValue
			case HDR_CONTENT_LENGTH_CMP:
				s.ContentLength = headerValue
			case HDR_CSEQ:
				s.Cseq = headerValue
			case HDR_FROM:
				s.From = headerValue
			case HDR_FROM_CMP:
				s.From = headerValue
			case HDR_MAX_FORWARDS:
				s.MaxForwards, _ = strconv.ParseInt(headerValue, 10, 8)
			case HDR_ORGANIZATION:
				s.Organization = headerValue
			case HDR_P_ASSERTED_IDENTITY:
				s.PAssertedIdVal = headerValue
			case HDR_PRIVACY:
				s.Privacy = headerValue
			case HDR_PROXY_AUTHENTICATE:
				s.ProxyAuthenticate = headerValue
			case HDR_RACK:
				s.Rack = headerValue
			case HDR_REASON:
				s.Reason = headerValue
			case HDR_RECORD_ROUTE:
				s.RecordRoute = append(s.RecordRoute, headerValue)
			case HDR_REMOTE_PARTY_ID:
				s.RemotePartyIdVal = headerValue
			case HDR_ROUTE:
				s.Route = append(s.RecordRoute, headerValue)
			case HDR_SERVER:
				s.Server = headerValue
			case HDR_SUPPORTED:
				s.Supported = append(s.RecordRoute, headerValue)
			case HDR_TO:
				s.To = headerValue
			case HDR_TO_CMP:
				s.To = headerValue
			case HDR_UNSUPPORTED:
				s.Unsupported = append(s.Unsupported, headerValue)
			case HDR_USER_AGENT:
				s.UserAgent = headerValue
			case HDR_VIA:
				s.Via = append(s.Via, headerValue)
			case HDR_VIA_CMP:
				s.Via = append(s.Via, headerValue)
			case HDR_WARNING:
				s.Warning = headerValue
			case HDR_WWW_AUTHENTICATE:
				s.WWWAuthenticate = headerValue
			default:
				s.OtherHeaders = append(s.OtherHeaders, headerValue)
			}
		}
	}
	return nil
}
