/**
* Homer Encapsulation Protocol v3
**/

package hep

import (
	"encoding/binary"
	"errors"
)

/*************************************
 Constants
*************************************/

// HEP ID
const HEP_ID = 0x48455033

// Generic Chunk Types
const (
	_ = iota // Don't want to assign zero here, but want to implicitly repeat this expression after...
	IP_PROTOCOL_FAMILY
	IP_PROTOCOL_ID
	IP4_SOURCE_ADDRESS
	IP4_DESTINATION_ADDRESS
	IP6_SOURCE_ADDRESS
	IP6_DESTINATION_ADDRESS
	SOURCE_PORT
	DESTINATION_PORT
	TIMESTAMP
	TIMESTAMP_MICRO
	PROTOCOL_TYPE // Maps to Protocol Types below
	CAPTURE_AGENT_ID
	KEEP_ALIVE_TIMER
	AUTHENTICATE_KEY
	PACKET_PAYLOAD
	COMPRESSED_PAYLOAD
	INTERNAL_C
)

var ProtocolFamilies []string
var Vendors []string
var ProtocolTypes []string

func init() {

	// Protocol Family Types - HEP3 Spec does not list these values out. Took IPv4 from an example.
	ProtocolFamilies = []string{
		"?",
		"?",
		"IPv4"}

	// Initialize vendors
	Vendors = []string{
		"None",
		"FreeSWITCH",
		"Kamailio",
		"OpenSIPS",
		"Asterisk",
		"Homer",
		"SipXecs",
	}

	// Initialize protocol types
	ProtocolTypes = []string{
		"Reserved",
		"SIP",
		"XMPP",
		"SDP",
		"RTP",
		"RTCP",
		"MGCP",
		"MEGACO",
		"M2UA",
		"M3UA",
		"IAX",
		"H322",
		"H321",
	}
}

// Define Struct for storing full HEP message
type HepMsg struct {
	IpProtocolFamily      int
	IpProtocolId          int
	Ip4SourceAddress      string
	Ip4DestinationAddress string
	Ip6SourceAddress      string
	Ip6DestinationAddress string
	SourcePort            int
	DestinationPort       int
	Timestamp             int32
	TimestampMicro        int32
	ProtocolType          int
	CaptureAgentId        int32
	KeepAliveTimer        int16
	AuthenticateKey       string
	Body                  string
}

func (hepMsg *HepMsg) Parse(udpPacket []byte) err {
	if binary.BigEndian.Uint32(buf[:3]) != HEP_ID {
		err := errors.New("Not a valid HEP3 packet - HEP3 ID is incorrect")
		return err
	}

	length := binary.BigEndian.Uint32(buf[4:5])
	currentByte := 6
	for currentByte < length {

	}

}

func (hepMsg *HepMsg) parseChunk(hepChunk []byte) (length, err) {
	chunkVendorId := binary.BigEndian.Uint16(hepChunk[:1])
	chunkType := binary.BigEndian.Uint16(hepChunk[2:3])
	chunkLength := binary.BigEndian.Uint32(hepChunk[4:5])
	chunkBody := hepChunk[6 : chunkLength-1]

	switch chunkType {
	case IP_PROTOCOL_FAMILY:
		hepMsg.IpProtocolFamily = binary.BigEndian.Uint16(chunkBody)
	case IP_PROTOCOL_ID:
		hepMsg.IpProtocolId = binary.BigEndian.Uint16(chunkBody)
	case IP4_SOURCE_ADDRESS:
		hepMsg.Ip4SourceAddress = binary.BigEndian.String(chunkBody)
	case IP4_DESTINATION_ADDRESS:
		hepMsg.Ip4DestinationAddress = binary.BigEndian.String(chunkBody)
	case IP6_SOURCE_ADDRESS:
		hepMsg.Ip6SourceAddress = binary.BigEndian.String(chunkBody)
	case IP6_DESTINATION_ADDRESS:
		hepMsg.Ip4DestinationAddress = binary.BigEndian.String(chunkBody)
	case SOURCE_PORT:
		hepMsg.SourcePort = binary.BigEndian.Uint16(chunkBody)
	case DESTINATION_PORT:
		hepMsg.DestinationPort = binary.BigEndian.Uint16(chunkBody)
	case TIMESTAMP:
		hepMsg.Timestamp = binary.BigEndian.Uint36(chunkBody)
	case TIMESTAMP_MICRO:
		hepMsg.TimestampMicro = binary.BigEndian.Uint36(chunkBody)
	case PROTOCOL_TYPE:
		hepMsg.ProtocolType = binary.BigEndian.Uint16(chunkBody)
	case CAPTURE_AGENT_ID:
	case KEEP_ALIVE_TIMER:
	case AUTHENTICATE_KEY:
	case PACKET_PAYLOAD:
	case COMPRESSED_PAYLOAD:
	case INTERNAL_C:
	}

}
