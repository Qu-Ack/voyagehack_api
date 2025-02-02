package mail

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Qu-Ack/voyagehack_api/services/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MailService struct {
	repo *MailRepo
}

func NewMailService(repo *MailRepo) *MailService {
	return &MailService{
		repo: repo,
	}
}

func (m *MailService) SendApplication(ctx context.Context, mail *Mail) (*Mail, error) {
	mail.Type = Application
	created, err := m.repo.createEmail(ctx, mail)

	if err != nil {
		return nil, err
	}

	err = m.repo.AddMailToSent(ctx, mail.Sender, mail.ID)
	if err != nil {
		return nil, err
	}

	err = m.repo.AddMailToReceived(ctx, mail.Receiver, mail.ID)
	if err != nil {
		return nil, err
	}

	return created, nil

}

func (m *MailService) InitializeMailBox(ctx context.Context, mailId string) (*MailBox, error) {
	return m.repo.createMailBox(ctx, mailId)
}

func (m *MailService) SendNormalMail(ctx context.Context, mail *Mail, requester user.PublicUser) (*Mail, error) {
	if requester.Role == "PATIENT" {
		return nil, errors.New("Can't send emails")
	}
	created, err := m.repo.createEmail(ctx, mail)

	if err != nil {
		return nil, err
	}

	err = m.repo.AddMailToSent(ctx, mail.Sender, mail.ID)
	if err != nil {
		return nil, err
	}

	err = m.repo.AddMailToReceived(ctx, mail.Receiver, mail.ID)
	if err != nil {
		return nil, err
	}
	fmt.Println(created)

	return created, nil
}

func (m *MailService) GetMailBox(ctx context.Context, userEmail string) (*PublicMailBox, error) {

	mailBox, err := m.repo.findMailBox(ctx, userEmail)

	if err != nil {
		return nil, err
	}
	fmt.Println(mailBox)

	return m.repo.AddMailsToMailBox(ctx, mailBox.ID)
}

func (m *MailService) ForwardMail(ctx context.Context, forwardTo string, mailId string, requester user.PublicUser) (*Mail, error) {

	if requester.Role == user.RolePatient {
		return nil, errors.New("not allowed to forward mails")
	}

	mailObjectId, err := primitive.ObjectIDFromHex(mailId)
	if err != nil {
		return nil, err
	}

	err = m.repo.addForwarded(ctx, forwardTo, mailObjectId)

	if err != nil {
		return nil, err
	}

	err = m.repo.AddMailToReceived(ctx, forwardTo, mailObjectId)
	if err != nil {
		return nil, err
	}

	err = m.repo.AddMailToSent(ctx, requester.Email, mailObjectId)
	if err != nil {
		return nil, err
	}

	return m.repo.getEmail(ctx, mailObjectId)

}

func (m *MailService) GetMail(ctx context.Context, mailId string, requester string) (*Mail, error) {
	mailBoxObjectId, err := primitive.ObjectIDFromHex(mailId)

	if err != nil {
		return nil, err
	}

	mail, err := m.repo.getEmail(ctx, mailBoxObjectId)

	if err != nil {
		return nil, err
	}

	if (requester == mail.Sender || requester == mail.Receiver) || slices.Contains(mail.ForwardedChain, requester) {
		return mail, nil
	}

	return nil, errors.New("You are not authorized to view this mail")

}
