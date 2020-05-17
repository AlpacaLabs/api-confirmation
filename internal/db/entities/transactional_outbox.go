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

type ConfirmEmailEvent struct {
	EventID        string
	Sent           bool
	EmailAddressID string
}

func NewConfirmEmailEvent(emailAddressID string) ConfirmEmailEvent {
	return ConfirmEmailEvent{
		EventID:        xid.New().String(),
		Sent:           false,
		EmailAddressID: emailAddressID,
	}
}

type ConfirmPhoneEvent struct {
	EventID       string
	Sent          bool
	PhoneNumberID string
}

func NewConfirmPhoneEvent(phoneNumberID string) ConfirmPhoneEvent {
	return ConfirmPhoneEvent{
		EventID:       xid.New().String(),
		Sent:          false,
		PhoneNumberID: phoneNumberID,
	}
}
