package entities

import (
	"github.com/golang/protobuf/proto"

	eventV1 "github.com/AlpacaLabs/protorepo-event-go/alpacalabs/event/v1"
	"github.com/rs/xid"
)

type Event struct {
	eventV1.EventInfo
	eventV1.TraceInfo
	Sent     bool
	Catalyst []byte
	Payload  []byte
}

func NewEvent(traceInfo eventV1.TraceInfo, catalyst, payload proto.Message) (Event, error) {
	var empty Event

	catalystBytes, err := proto.Marshal(catalyst)
	if err != nil {
		return empty, err
	}

	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return empty, err
	}

	return Event{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo: traceInfo,
		Sent:      false,
		Catalyst:  catalystBytes,
		Payload:   payloadBytes,
	}, nil
}
