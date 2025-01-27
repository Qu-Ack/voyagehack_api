package graph

import (
	"github.com/Qu-Ack/voyagehack_api/services/mail"
	"github.com/Qu-Ack/voyagehack_api/services/messaging"
	"github.com/Qu-Ack/voyagehack_api/services/observers"
	"github.com/Qu-Ack/voyagehack_api/services/upload"
	"github.com/Qu-Ack/voyagehack_api/services/user"
)

// This file will not be regenerated automatically.
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService      *user.UserService
	ObserverService  *observers.ObserverService
	MailService      *mail.MailService
	MessagingService *messaging.MessagingService
	UploadService    *upload.UploadService
}
