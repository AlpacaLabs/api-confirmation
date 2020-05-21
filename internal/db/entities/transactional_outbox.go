package entities

import (
	"context"

	"github.com/golang/protobuf/proto"

	"github.com/AlpacaLabs/go-kontext"

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

func NewEvent(ctx context.Context, catalyst, payload proto.Message) (Event, error) {
	var empty Event

	traceInfo := kontext.GetTraceInfo(ctx)

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
