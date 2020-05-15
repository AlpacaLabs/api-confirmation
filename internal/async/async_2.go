package async

// TODO
// 1) read event from transactional outbox
// 2) write to send-email-request or send-sms-request topic
// 3) mark txob record as sent

// TODO does every pod read from txob?
//  wouldn't that result in duplicate messages sent to Hermes?
