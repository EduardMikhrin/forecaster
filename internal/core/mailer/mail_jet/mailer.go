package mail_jet

import (
	"context"
	"fmt"
	"github.com/EduardMikhrin/forecaster/internal"
	"github.com/EduardMikhrin/forecaster/internal/core/mailer"
	"github.com/fatih/structs"
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
)

type mailJet struct {
	client             *mailjet.Client
	from               mailjet.RecipientV31
	verificationTmplId int64
	infoTmplId         int64

	ctx context.Context
	log *logan.Entry
}

func (s *mailJet) SendVerificationEmail(to string, payload interface{}) error {
	code, ok := payload.(string)
	if !ok {
		return errors.New("payload is not of type string")
	}

	var vars = map[string]interface{}{
		"code": code,
	}

	if err := s.sendList([]string{to}, mailer.VerificationSubject, s.verificationTmplId, vars); err != nil {
		return errors.Wrap(err, "failed to send verification email")
	}

	return nil
}

func (s *mailJet) SendInfoEmail(to []string, payload interface{}) error {
	info, ok := payload.(*internal.WeatherPayload)
	if !ok {
		return errors.New("payload is not of type WeatherPayload")
	}

	if err := s.sendList(to, mailer.InfoSubject, s.infoTmplId, structs.Map(info)); err != nil {
		return errors.Wrap(err, "failed to send info email")
	}

	return nil
}

func (m *mailJet) sendList(
	emails []string,
	subject string,
	templateID int64,
	variables map[string]interface{},
) error {
	select {
	case <-m.ctx.Done():
		m.log.Debug("context closed")
		return nil
	default:
	}

	var recipients mailjet.RecipientsV31
	for _, email := range emails {
		recipients = append(recipients, mailjet.RecipientV31{
			Email: email,
		})
	}

	messages := mailjet.MessagesV31{
		Info: []mailjet.InfoMessagesV31{
			{
				From:             &m.from,
				To:               &recipients,
				TemplateID:       int(templateID),
				TemplateLanguage: true,
				Subject:          subject,
				Variables:        variables,
			},
		},
	}

	res, err := m.client.SendMailV31(&messages)
	if err != nil {
		return errors.Wrap(err, "failed to send mail")
	}

	if len(res.ResultsV31) == 0 {
		return fmt.Errorf("mailjet: message not sent")
	}

	return nil
}

func NewNotifier(client *mailjet.Client, from mailjet.RecipientV31, verificationTmplId, infoTmplId int64,
	ctx context.Context, log *logan.Entry) mailer.Mailer {
	return &mailJet{
		client:             client,
		from:               from,
		verificationTmplId: verificationTmplId,
		infoTmplId:         infoTmplId,
		ctx:                ctx,
		log:                log,
	}
}
