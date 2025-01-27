package messaging

import (
	"context"
	"errors"
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

	participantObjectId, err := primitive.ObjectIDFromHex(participantId)
	if err != nil {
		return nil, err
	}

	requestUserObjectId, err := primitive.ObjectIDFromHex(requester.ID)

	if err != nil {
		return nil, err
	}

	var particpants []primitive.ObjectID

	particpants = append(particpants, participantObjectId, requestUserObjectId)

	return m.repo.createRoom(ctx, &Room{
		Participants: particpants,
		Messages:     make([]Message, 0),
	})

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
		return m.repo.addMessage(ctx, roomObjectId, message)
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

	if slices.Contains(room.Participants, userObjectId) {
		return room, nil
	} else {
		return nil, errors.New("Participant should be inside the room")
	}

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
