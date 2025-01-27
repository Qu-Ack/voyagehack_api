package mail

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MailBox struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Email    string               `bson:"email,omitempty" json:"email"`
	Sent     []primitive.ObjectID `bson:"sent,omitempty" json:"sent"`
	Received []primitive.ObjectID `bson:"received,omitempty" json:"received"`
}

type PublicMailBox struct {
	MailBox        MailBox
	SentEMails     []Mail
	ReceivedEmails []Mail
}

type EMAIL string

const (
	Application EMAIL = "APPLICATION"
	Normal      EMAIL = "NORMAL"
)

type Mail struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Sender    string             `bson:"sender,omitempty" json:"sender"`
	Receiver  string             `bson:"receiver,omitempty" json:"receiver"`
	Content   string             `bson:"content,omitempty" json:"content"`
	Documents []*string          `bson:"documents,omitempty" json:"documents"`
	Type      EMAIL              `bson:"type,omitempty" json:"type"`
	CreatedAt primitive.DateTime `bson:"createdAt,omitempty" json:"createdAt"`
}
