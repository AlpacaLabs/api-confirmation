package async

const (
	topicNameForCreatingPhoneNumberConfirmationCode = "create-phone-number-confirmation-code-request"
)

// Async 1
// read event from topic
// have we seen this event before? check event's ID in DB
// if not, persist event to transactional outbox
//
// Async 2
// read event from transactional outbox
// send email / sms
// update event entity as sent
