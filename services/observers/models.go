package observers

import "go.mongodb.org/mongo-driver/bson/primitive"

type MailBoxSubscriptionResopnse struct {
	Sent     Mail `json:"sent"`
	Received Mail `json:"received"`
}

type MessageSubscriptionResponse struct {
	RoomId  string  `json:"roomid"`
	Message Message `json:"message"`
}

type Message struct {
	Author  string `json:"author"`
	Content string `json:"content"`
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
