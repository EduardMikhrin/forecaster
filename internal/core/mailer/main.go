package mailer

type Mailer interface {
	SendVerificationEmail(to string, payload interface{}) error
	SendInfoEmail(to []string, payload interface{}) error
}
