package mail

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MailRepo struct {
	db *mongo.Database
}

func NewMailRepo(db *mongo.Database) *MailRepo {
	return &MailRepo{
		db: db,
	}
}

func (m *MailRepo) createEmail(ctx context.Context, mail *Mail) (*Mail, error) {
	mailCollection := m.db.Collection("mails")

	insertOneResult, err := mailCollection.InsertOne(ctx, mail)

	if err != nil {
		return nil, err
	}

	mail.ID = insertOneResult.InsertedID.(primitive.ObjectID)
	return mail, nil

}

func (m *MailRepo) getEmail(ctx context.Context, mailId primitive.ObjectID) (*Mail, error) {
	mailCollection := m.db.Collection("mails")

	var mail Mail
	err := mailCollection.FindOne(ctx, bson.M{"_id": mailId}).Decode(&mail)

	if err != nil {
		return nil, err
	}

	return &mail, nil
}

func (m *MailRepo) findMailBox(ctx context.Context, emailId string) (*MailBox, error) {
	mailBoxCollection := m.db.Collection("mailbox")

	var mailBox MailBox
	err := mailBoxCollection.FindOne(ctx, bson.M{"email": emailId}).Decode(&mailBox)

	if err != nil {
		return nil, err
	}

	return &mailBox, nil
}

func (m *MailRepo) AddMailToSent(ctx context.Context, mailId string, mail primitive.ObjectID) error {
	mailBoxCollection := m.db.Collection("mailbox")

	result := mailBoxCollection.FindOneAndUpdate(ctx, bson.M{"Email": mailId}, bson.D{{"$push", bson.D{{"sent", mail}}}})

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (m *MailRepo) AddMailToReceived(ctx context.Context, mailId string, mail primitive.ObjectID) error {
	mailBoxCollection := m.db.Collection("mailbox")

	result := mailBoxCollection.FindOneAndUpdate(ctx, bson.M{"Email": mailId}, bson.D{{"$push", bson.D{{"received", mail}}}})

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (m *MailRepo) AddMailsToMailBox(ctx context.Context, mailBoxId primitive.ObjectID) (*PublicMailBox, error) {
	mailBoxCollection := m.db.Collection("mailbox")
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"_id": mailBoxId}}},
		{{"$lookup", bson.M{
			"from":         "mails",
			"localField":   "sent",
			"foreignField": "_id",
			"as":           "SentEmails",
		}}},
		{{"$lookup", bson.M{
			"from":         "mails",
			"localField":   "sent",
			"foreignField": "_id",
			"as":           "ReceivedEmails",
		}}},
		{{"$project", bson.M{
			"sent":           1,
			"received":       1,
			"Email":          1,
			"SentEmails":     1,
			"ReceivedEmails": 1,
		}}},
	}

	cursor, err := mailBoxCollection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var publicMailBox PublicMailBox
	if cursor.Next(ctx) {
		err := cursor.Decode(&publicMailBox)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("mail box not found")
	}

	return &publicMailBox, nil
}
