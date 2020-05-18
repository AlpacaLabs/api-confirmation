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

func NewSendEmailEvent(ctx context.Context, codeID string) SendEmailEvent {
	return SendEmailEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo: eventV1.TraceInfo{
			TraceId: kontext.GetTraceID(ctx),
			Sampled: false,
		},
		Sent:   false,
		CodeID: codeID,
	}
}

type SendPhoneEvent struct {
	eventV1.EventInfo
	eventV1.TraceInfo
	Sent   bool
	CodeID string
}

func NewSendPhoneEvent(ctx context.Context, codeID string) SendPhoneEvent {
	return SendPhoneEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo: eventV1.TraceInfo{
			TraceId: kontext.GetTraceID(ctx),
			Sampled: false,
		},
		Sent:   false,
		CodeID: codeID,
	}
}

type ConfirmEmailEvent struct {
	eventV1.EventInfo
	eventV1.TraceInfo
	Sent           bool
	EmailAddressID string
}

func NewConfirmEmailEvent(ctx context.Context, emailAddressID string) ConfirmEmailEvent {
	return ConfirmEmailEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo: eventV1.TraceInfo{
			TraceId: kontext.GetTraceID(ctx),
			Sampled: false,
		},
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
	return ConfirmPhoneEvent{
		EventInfo: eventV1.EventInfo{
			EventId: xid.New().String(),
		},
		TraceInfo: eventV1.TraceInfo{
			TraceId: kontext.GetTraceID(ctx),
			Sampled: false,
		},
		Sent:          false,
		PhoneNumberID: phoneNumberID,
	}
}
