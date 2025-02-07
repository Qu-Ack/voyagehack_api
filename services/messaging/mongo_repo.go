package messaging

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepo struct {
	db *mongo.Database
}

func NewMessageRepo(db *mongo.Database) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

func (m *MessageRepo) createRoom(ctx context.Context, room *Room) (*Room, error) {
	fmt.Println("in create room")
	chatCollection := m.db.Collection("chat")

	fmt.Println("before insertion")
	insertResult, err := chatCollection.InsertOne(ctx, room)

	if err != nil {
		return nil, err
	}

	fmt.Println("after insertion")

	room.ID = insertResult.InsertedID.(primitive.ObjectID)

	return room, nil
}

func (m *MessageRepo) closeRoom(ctx context.Context, roomId primitive.ObjectID) (*Room, error) {
	chatCollection := m.db.Collection("chat")

	var room Room
	err := chatCollection.FindOneAndUpdate(ctx, bson.M{"_id": roomId}, bson.D{{"$set", bson.D{{"state", "CLOSE"}}}}).Decode(&room)

	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (m *MessageRepo) openRoom(ctx context.Context, roomId primitive.ObjectID) (*Room, error) {

	chatCollection := m.db.Collection("chat")

	var room Room
	err := chatCollection.FindOneAndUpdate(ctx, bson.M{"_id": roomId}, bson.D{{"$set", bson.D{{"state", "OPEN"}}}}).Decode(&room)

	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (m *MessageRepo) addMessage(ctx context.Context, roomId primitive.ObjectID, message *Message) (*Room, error) {
	chatCollection := m.db.Collection("chat")

	var room Room
	err := chatCollection.FindOneAndUpdate(ctx, bson.M{"_id": roomId}, bson.D{{"$push", bson.D{{"messages", message}}}}).Decode(&room)
	fmt.Println(room)

	if err != nil {
		return nil, err
	}

	return &room, nil

}

func (m *MessageRepo) getRoomsById(ctx context.Context, resourceId primitive.ObjectID) ([]*Room, error) {
	chatCollection := m.db.Collection("chat")

	cursor, err := chatCollection.Find(ctx, bson.M{"participants": resourceId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []*Room
	for cursor.Next(ctx) {
		var room Room
		if err := cursor.Decode(&room); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	// Check for errors after iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (m *MessageRepo) getRoom(ctx context.Context, roomId primitive.ObjectID) (*Room, error) {
	chatCollection := m.db.Collection("chat")

	var room Room
	err := chatCollection.FindOne(ctx, bson.M{"_id": roomId}).Decode(&room)

	if err != nil {
		return nil, err
	}

	return &room, nil
}
