package entities

import "github.com/rs/xid"

type SendEmailEvent struct {
	EventID string
	Sent    bool
	CodeID  string
}

func NewSendEmailEvent(codeID string) SendEmailEvent {
	return SendEmailEvent{
		EventID: xid.New().String(),
		Sent:    false,
		CodeID:  codeID,
	}
}

type SendPhoneEvent struct {
	EventID string
	Sent    bool
	CodeID  string
}

func NewSendPhoneEvent(codeID string) SendPhoneEvent {
	return SendPhoneEvent{
		EventID: xid.New().String(),
		Sent:    false,
		CodeID:  codeID,
	}
}
