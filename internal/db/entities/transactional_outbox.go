package entities

import (
	"context"

	"github.com/AlpacaLabs/go-kontext"

	eventV1 "github.com/AlpacaLabs/protorepo-event-go/alpacalabs/event/v1"
	"github.com/rs/xid"
)

type SendEmailEvent struct {
	eventV1.EventInfo
	eventV1.TraceInfo
	Sent   bool
	CodeID string
}

func NewSendEmailEvent(traceInfo eventV1.TraceInfo, codeID string) SendEmailEvent {
	return SendEmailEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo: traceInfo,
		Sent:      false,
		CodeID:    codeID,
	}
}

type SendPhoneEvent struct {
	eventV1.EventInfo
	eventV1.TraceInfo
	Sent   bool
	CodeID string
}

func NewSendPhoneEvent(traceInfo eventV1.TraceInfo, codeID string) SendPhoneEvent {
	return SendPhoneEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo: traceInfo,
		Sent:      false,
		CodeID:    codeID,
	}
}

type ConfirmEmailEvent struct {
	eventV1.EventInfo
	eventV1.TraceInfo
	Sent           bool
	EmailAddressID string
}

func NewConfirmEmailEvent(ctx context.Context, emailAddressID string) ConfirmEmailEvent {
	traceInfo := kontext.GetTraceInfo(ctx)
	return ConfirmEmailEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo:      traceInfo,
		Sent:           false,
		EmailAddressID: emailAddressID,
	}
}

type ConfirmPhoneEvent struct {
	eventV1.EventInfo
	eventV1.TraceInfo
	Sent          bool
	PhoneNumberID string
}

func NewConfirmPhoneEvent(ctx context.Context, phoneNumberID string) ConfirmPhoneEvent {
	traceInfo := kontext.GetTraceInfo(ctx)
	return ConfirmPhoneEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo:     traceInfo,
		Sent:          false,
		PhoneNumberID: phoneNumberID,
	}
}
