package mailer

type INotifier interface {
	SendVerificationEmail(to string, payload interface{}) error
	SendInfoEmail(to []string, payload interface{}) error
}
