package entities

import (
	"time"

	clock "github.com/AlpacaLabs/go-timestamp"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

type PhoneNumberConfirmationCode struct {
	ID            string
	PhoneNumberID string
	Code          string
	CreatedAt     time.Time
	ExpiresAt     time.Time
	Used          bool
	Stale         bool
}

func NewPhoneNumberConfirmationCodeFromProtobuf(c confirmationV1.PhoneNumberConfirmationCode) PhoneNumberConfirmationCode {
	return PhoneNumberConfirmationCode{
		ID:            c.Id,
		PhoneNumberID: c.PhoneNumberId,
		Code:          c.Code,
		CreatedAt:     clock.TimestampToTime(c.CreatedAt),
		ExpiresAt:     clock.TimestampToTime(c.ExpiresAt),
		Used:          c.Used,
		Stale:         c.Stale,
	}
}

func (c PhoneNumberConfirmationCode) ToProtobuf() *confirmationV1.PhoneNumberConfirmationCode {
	return &confirmationV1.PhoneNumberConfirmationCode{
		Id:            c.ID,
		PhoneNumberId: c.PhoneNumberID,
		Code:          c.Code,
		CreatedAt:     clock.TimeToTimestamp(c.CreatedAt),
		ExpiresAt:     clock.TimeToTimestamp(c.ExpiresAt),
		Used:          c.Used,
		Stale:         c.Stale,
	}
}

type EmailAddressConfirmationCode struct {
	ID             string
	EmailAddressID string
	Code           string
	CreatedAt      time.Time
	ExpiresAt      time.Time
	Used           bool
	Stale          bool
}

func NewEmailAddressConfirmationCodeFromProtobuf(c confirmationV1.EmailAddressConfirmationCode) EmailAddressConfirmationCode {
	return EmailAddressConfirmationCode{
		ID:             c.Id,
		EmailAddressID: c.EmailAddressId,
		Code:           c.Code,
		CreatedAt:      clock.TimestampToTime(c.CreatedAt),
		ExpiresAt:      clock.TimestampToTime(c.ExpiresAt),
		Used:           c.Used,
		Stale:          c.Stale,
	}
}

func (c EmailAddressConfirmationCode) ToProtobuf() *confirmationV1.EmailAddressConfirmationCode {
	return &confirmationV1.EmailAddressConfirmationCode{
		Id:             c.ID,
		EmailAddressId: c.EmailAddressID,
		Code:           c.Code,
		CreatedAt:      clock.TimeToTimestamp(c.CreatedAt),
		ExpiresAt:      clock.TimeToTimestamp(c.ExpiresAt),
		Used:           c.Used,
		Stale:          c.Stale,
	}
}
