package messaging

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/Qu-Ack/voyagehack_api/services/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessagingService struct {
	repo *MessageRepo
}

func NewMessageService(repo *MessageRepo) *MessagingService {
	return &MessagingService{
		repo: repo,
	}
}

func (m *MessagingService) CreateRoom(ctx context.Context, participantId string, requester user.PublicUser) (*Room, error) {
	if requester.Role != user.RoleDoctor {
		return nil, errors.New("Only Doctors Can Create A chat room")
	}
	fmt.Println("role check")

	participantObjectId, err := primitive.ObjectIDFromHex(participantId)
	if err != nil {
		return nil, err
	}
	fmt.Println(participantObjectId)

	requestUserObjectId, err := primitive.ObjectIDFromHex(requester.ID)

	if err != nil {
		return nil, err
	}
	fmt.Println(requestUserObjectId)

	participants := make([]primitive.ObjectID, 0)

	participants = append(participants, participantObjectId, requestUserObjectId)
	fmt.Println(participants)
	messages := make([]Message, 0)
	room := Room{
		Participants: participants,
		Messages:     messages,
		State:        "OPEN",
	}

	returnedRoom, err := m.repo.createRoom(ctx, &room)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("room created success")

	return returnedRoom, nil

}

func (m *MessagingService) SendMessage(ctx context.Context, roomId string, requester user.PublicUser, message *Message) (*Room, error) {

	roomObjectId, err := primitive.ObjectIDFromHex(roomId)

	if err != nil {
		return nil, err
	}

	room, err := m.repo.getRoom(ctx, roomObjectId)

	if err != nil {
		return nil, err
	}

	if room.State == "CLOSE" {
		return nil, errors.New("room should be open")
	}

	userObjectId, err := primitive.ObjectIDFromHex(requester.ID)

	if err != nil {
		return nil, err
	}

	if slices.Contains(room.Participants, userObjectId) {
		_, err := m.repo.addMessage(ctx, roomObjectId, message)
		if err != nil {
			return nil, err
		}

		room, err = m.repo.getRoom(ctx, roomObjectId)
		if err != nil {
			return nil, err
		}
		return room, nil
	} else {
		return nil, errors.New("requester should be part of the room")
	}

}

func (m *MessagingService) GetRoom(ctx context.Context, roomId string, requester user.PublicUser) (*Room, error) {
	roomObjectId, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return nil, err
	}
	userObjectId, err := primitive.ObjectIDFromHex(requester.ID)
	if err != nil {
		return nil, err
	}

	room, err := m.repo.getRoom(ctx, roomObjectId)
	if err != nil {
		return nil, err
	}

	if room.State == "CLOSE" {
		return nil, errors.New("room is closed")
	}

	if slices.Contains(room.Participants, userObjectId) {
		return room, nil
	} else {
		return nil, errors.New("Participant should be inside the room")
	}

}

func (m *MessagingService) GetRoomsByID(ctx context.Context, userId string) ([]*Room, error) {

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	rooms, err := m.repo.getRoomsById(ctx, userObjectId)

	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (m *MessagingService) ChangeRoomState(ctx context.Context, state string, requester user.PublicUser, roomId string) (*Room, error) {

	roomObjectId, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return nil, err
	}

	room, err := m.repo.getRoom(ctx, roomObjectId)
	if err != nil {
		return nil, err
	}

	switch state {
	case "OPEN":
		if requester.Role == user.RoleDoctor {
			return m.repo.openRoom(ctx, room.ID)
		} else {
			return nil, errors.New("Only doctor can open a room")
		}
	case "CLOSE":
		return m.repo.closeRoom(ctx, roomObjectId)
	default:
		return nil, errors.New("choose a valid state")
	}

}
