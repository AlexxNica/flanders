package sip

import (
	"testing"
)

func TestBasicStartLineParse(t *testing.T) {
	t.Log("Test basic start line parse")
	startLine, err := parseStartLine("SIP/2.0 180 Ringing")
	if err != nil {
		t.Fatalf("StartLine.Parse: Error parsing basic start line: %s", err)
	}

	t.Log("Parsing basic 180-Ringing message")

	if startLine.Body != "SIP/2.0 180 Ringing" {
		t.Fatalf("StartLine.Parse: Expected body of sip message to be different: %s", startLine.Body)
	}

	if startLine.Type != SIP_RESPONSE {
		t.Fatalf("StartLine.Parse: Expected type of start line to be response but was: %s", startLine.Type)
	}
}
