package mail

import (
	"context"
	"errors"
	"fmt"
	"slices"

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
	fmt.Println("in create email")
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

func (m *MailRepo) addForwarded(ctx context.Context, forwardedTo string, mailId primitive.ObjectID) error {
	mailBoxCollection := m.db.Collection("mails")

	result := mailBoxCollection.FindOneAndUpdate(ctx, bson.M{"_id": mailId}, bson.D{{"$push", bson.D{{"forwardedChain", forwardedTo}}}})

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (m *MailRepo) AddMailToSent(ctx context.Context, mailId string, mail primitive.ObjectID) error {
	fmt.Println("add mails to sent")
	mailBoxCollection := m.db.Collection("mailbox")

	result := mailBoxCollection.FindOneAndUpdate(ctx, bson.M{"email": mailId}, bson.D{{"$push", bson.D{{"sent", mail}}}})

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (m *MailRepo) AddMailToReceived(ctx context.Context, mailId string, mail primitive.ObjectID) error {
	fmt.Println("add mails to received")
	mailBoxCollection := m.db.Collection("mailbox")

	result := mailBoxCollection.FindOneAndUpdate(ctx, bson.M{"email": mailId}, bson.D{{"$push", bson.D{{"received", mail}}}})

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (m *MailRepo) checkUserInMail(ctx context.Context, emailId string, mailId primitive.ObjectID) error {
	mailBoxCollection := m.db.Collection("mailbox")

	var mail Mail
	err := mailBoxCollection.FindOne(ctx, bson.M{"_id": mailId}).Decode(&mail)

	if err != nil {
		return err
	}

	if mail.Sender == string(emailId) || mail.Receiver == string(emailId) {
		return nil
	} else if slices.Contains(mail.ForwardedChain, emailId) {
		return nil
	} else {
		return errors.New("can't forward mail")
	}

}

func (m *MailRepo) createMailBox(ctx context.Context, mailId string) (*MailBox, error) {
	mailBoxCollection := m.db.Collection("mailbox")

	mailbox := MailBox{
		Email: mailId,
	}

	_, err := mailBoxCollection.InsertOne(ctx, mailbox)

	if err != nil {
		return nil, err
	}

	return &mailbox, nil
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
			"localField":   "received", // Changed from "sent" to "received"
			"foreignField": "_id",
			"as":           "ReceivedEmails",
		}}},
		{{"$project", bson.M{
			"_id":            1,
			"email":          1,
			"sent":           1,
			"received":       1,
			"SentEmails":     1,
			"ReceivedEmails": 1,
		}}},
	}

	cursor, err := mailBoxCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result struct {
		ID             primitive.ObjectID   `bson:"_id"`
		Email          string               `bson:"email"`
		Sent           []primitive.ObjectID `bson:"sent"`
		Received       []primitive.ObjectID `bson:"received"`
		SentEmails     []Mail               `bson:"SentEmails"`
		ReceivedEmails []Mail               `bson:"ReceivedEmails"`
	}

	if cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}

		publicMailBox := PublicMailBox{
			MailBox: MailBox{
				ID:       result.ID,
				Email:    result.Email,
				Sent:     result.Sent,
				Received: result.Received,
			},
			SentEMails:     result.SentEmails,
			ReceivedEmails: result.ReceivedEmails,
		}
		return &publicMailBox, nil
	}

	return nil, errors.New("mail box not found")
}
