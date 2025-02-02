package graph

import (
	"fmt"
	"time"

	"github.com/Qu-Ack/voyagehack_api/api/graph/model"
	"github.com/Qu-Ack/voyagehack_api/services/hospital"
	"github.com/Qu-Ack/voyagehack_api/services/mail"
	"github.com/Qu-Ack/voyagehack_api/services/messaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/rand"
)

type AuthenticatedUser struct {
	ID    string
	Email string
	Role  string
}

type contextKey string

const UserContextKey contextKey = "user"

func convertObjectIDToStringSlice(ids []primitive.ObjectID) []*string {
	result := make([]*string, len(ids))
	for i, id := range ids {
		str := id.Hex()
		result[i] = &str
	}
	return result
}
func convertHospitalReviewsToModelReviews(reviews []hospital.Review) []*model.Review {
	result := make([]*model.Review, len(reviews))

	for i, review := range reviews {
		newReview := &model.Review{
			Content: review.Content,
			Author:  review.Author.Hex(),
		}
		result[i] = newReview
	}

	return result
}

func convertToModelMail(m mail.Mail) *model.Mail {
	return &model.Mail{
		ID:        m.ID.Hex(),
		Sender:    m.Sender,
		Receiver:  m.Receiver,
		Documents: m.Documents,
		Content:   m.Content,
		Type:      model.EmailType(m.Type),
		CreatedAt: m.CreatedAt.Time().String(),
	}
}
func convertMailSlice(mails []mail.Mail) []*model.Mail {
	result := make([]*model.Mail, len(mails))
	for i, m := range mails {
		result[i] = convertToModelMail(m)
	}
	return result
}

func convertMessages(messages []messaging.Message) []*model.Message {
	if messages == nil {
		return nil
	}
	result := make([]*model.Message, len(messages))
	for i, msg := range messages {
		result[i] = &model.Message{
			Content: msg.Content,
			Author:  msg.Author,
		}
	}
	return result
}

func convertParticipants(participants []primitive.ObjectID) []string {
	if participants == nil {
		return nil
	}
	result := make([]string, len(participants))
	for i, participant := range participants {
		result[i] = participant.Hex() // Convert ObjectID to string
	}
	return result
}

func selectRandomParticipant(participants []primitive.ObjectID) (primitive.ObjectID, error) {
	if len(participants) == 0 {
		return primitive.NilObjectID, fmt.Errorf("no participants available")
	}

	// Seed should only be called once in init()
	rand.Seed(uint64(time.Now().UnixNano()))

	randomIndex := rand.Intn(len(participants))

	return participants[randomIndex], nil
}

func convertStringToPointerStringArray(arr []string) []*string {
	if arr == nil {
		return nil
	}

	result := make([]*string, len(arr))

	for i, entity := range arr {
		result[i] = &entity
	}
	return result
}
