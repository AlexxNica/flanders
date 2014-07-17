package gosip

import (
	"testing"
)

func TestParse(t *testing.T) {
	t.Log("Test basic parse")
	msg := &SipMsg{}

	t.Log("Parsing basic 180-Ringing message")

	msg.Parse(`SIP/2.0 180 Ringing
Via: SIP/2.0/UDP 10.25.0.10:5060;branch=z9hG4bK7937.0761dd14.0
Via: SIP/2.0/UDP 10.25.2.32;received=10.25.2.32;rport=5060;branch=z9hG4bK3apF9tNpS9evp
From: "WIRELESS CALLER" <sip:4326612233@10.25.2.32>;tag=HKN3D9BKrQvcp
To: <sip:9725915920@10.25.0.10>;tag=gK0ad93ef2
Call-ID: 286b0616-766b-1232-a6b9-00259078f838
CSeq: 61480112 INVITE
Record-Route: <sip:10.25.0.10:5060;lr;ftag=HKN3D9BKrQvcp>
Contact: <sip:9725915920@66.2.204.94:5060>
Allow: INVITE,ACK,CANCEL,BYE,REGISTER,REFER,INFO,SUBSCRIBE,NOTIFY,PRACK,UPDATE,OPTIONS,MESSAGE,PUBLISH
P-Asserted-Identity: ". ." <sip:9725915920@10.25.2.32;user=phone>
Content-Length: 0

`)
	if msg.StartLine != "SIP/2.0 180 Ringing" {
		t.Fatalf("SipMsg.Parse: Expected first line of sip message to be different: %s", msg.StartLine)
	}

	if msg.Via[0] != "SIP/2.0/UDP 10.25.0.10:5060;branch=z9hG4bK7937.0761dd14.0" {
		t.Fatalf("SipMsg.Parse: Expected 2 values for Via property, but first message was wrong: %s", msg.Via[0])
	}
}
