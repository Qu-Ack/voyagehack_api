package messaging

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Messages     []Message            `bson:"messages,omitempty"`
	Participants []primitive.ObjectID `bson:"participants,omitempty"`
	State        string               `bson:"state,omitempty"`
}

type Message struct {
	Content string `bson:"content,omitempty"`
	Author  string `bson:"author,omitempty"`
}
