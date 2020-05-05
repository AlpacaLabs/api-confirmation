package entities

import (
	"time"

	clock "github.com/AlpacaLabs/go-timestamp"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/guregu/null"
)

type PhoneNumberConfirmationCode struct {
	ID            string
	PhoneNumberID string
	Code          string
	CreatedAt     null.Time
	ExpiresAt     null.Time
	Used          bool
	Stale         bool
}

func NewPhoneNumberConfirmationCodeFromProtobuf(c confirmationV1.PhoneNumberConfirmationCode) PhoneNumberConfirmationCode {
	return PhoneNumberConfirmationCode{
		ID:            c.Id,
		PhoneNumberID: c.PhoneNumberId,
		Code:          c.Code,
		CreatedAt:     timestampToNullTime(c.CreatedAt),
		ExpiresAt:     timestampToNullTime(c.ExpiresAt),
		Used:          c.Used,
		Stale:         c.Stale,
	}
}

func (c PhoneNumberConfirmationCode) ToProtobuf() *confirmationV1.PhoneNumberConfirmationCode {
	return &confirmationV1.PhoneNumberConfirmationCode{
		Id:            c.ID,
		PhoneNumberId: c.PhoneNumberID,
		Code:          c.Code,
		CreatedAt:     clock.TimeToTimestamp(c.CreatedAt.ValueOrZero()),
		ExpiresAt:     clock.TimeToTimestamp(c.ExpiresAt.ValueOrZero()),
		Used:          c.Used,
		Stale:         c.Stale,
	}
}

type EmailAddressConfirmationCode struct {
	ID             string
	EmailAddressID string
	Code           string
	CreatedAt      null.Time
	ExpiresAt      null.Time
	Used           bool
	Stale          bool
}

func NewEmailAddressConfirmationCodeFromProtobuf(c confirmationV1.EmailAddressConfirmationCode) EmailAddressConfirmationCode {
	return EmailAddressConfirmationCode{
		ID:             c.Id,
		EmailAddressID: c.EmailAddressId,
		Code:           c.Code,
		CreatedAt:      timestampToNullTime(c.CreatedAt),
		ExpiresAt:      timestampToNullTime(c.ExpiresAt),
		Used:           c.Used,
		Stale:          c.Stale,
	}
}

func (c EmailAddressConfirmationCode) ToProtobuf() *confirmationV1.EmailAddressConfirmationCode {
	return &confirmationV1.EmailAddressConfirmationCode{
		Id:             c.ID,
		EmailAddressId: c.EmailAddressID,
		Code:           c.Code,
		CreatedAt:      clock.TimeToTimestamp(c.CreatedAt.ValueOrZero()),
		ExpiresAt:      clock.TimeToTimestamp(c.ExpiresAt.ValueOrZero()),
		Used:           c.Used,
		Stale:          c.Stale,
	}
}

// TODO add to AlpacaLabs/go-timestamp?
func timestampToNullTime(in *timestamp.Timestamp) null.Time {
	t := clock.TimestampToTime(in)

	var nt null.Time
	if t.IsZero() {
		nt = null.NewTime(time.Time{}, false)
	} else {
		nt = null.TimeFrom(t)
	}
	return nt
}
